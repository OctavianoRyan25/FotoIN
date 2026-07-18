package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
)

type OrderUsecase struct {
	orderRepo      domain.OrderRepository
	galleryRepo    domain.GalleryRepository
	eventRepo      domain.EventRepository
	photoRepo      domain.PhotoRepository
	userRepo       domain.UserRepository
	walletRepo     domain.WalletRepository
	s3             *infra.S3Client
	mail           *infra.MailClient
	txManager      domain.TxManager
	paymentService domain.PaymentService
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	galleryRepo domain.GalleryRepository,
	eventRepo domain.EventRepository,
	photoRepo domain.PhotoRepository,
	userRepo domain.UserRepository,
	walletRepo domain.WalletRepository,
	s3 *infra.S3Client,
	mail *infra.MailClient,
	txManager domain.TxManager,
	paymentService domain.PaymentService,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:      orderRepo,
		galleryRepo:    galleryRepo,
		eventRepo:      eventRepo,
		photoRepo:      photoRepo,
		userRepo:       userRepo,
		walletRepo:     walletRepo,
		s3:             s3,
		mail:           mail,
		txManager:      txManager,
		paymentService: paymentService,
	}
}

// ─── Input / Output ─────────────────────────────────────────────────────────

type CreateOrderInput struct {
	BuyerEmail string   `json:"buyer_email" validate:"required,email"`
	GalleryID  string   `json:"gallery_id"` // one of gallery_id or event_id required
	EventID    string   `json:"event_id"`
	PhotoIDs   []string `json:"photo_ids" validate:"required,min=1"`
}

type OrderSummary struct {
	Order         *domain.Order                 `json:"order"`
	QRCodePayment *domain.QRCodePaymentResponse `json:"qr_code_payment,omitempty"`
}

// ─── CreateOrder ─────────────────────────────────────────────────────────────

func (u *OrderUsecase) CreateOrder(ctx context.Context, input CreateOrderInput) (*OrderSummary, error) {
	if input.GalleryID == "" && input.EventID == "" {
		return nil, errors.New("gallery_id or event_id is required")
	}

	// Parse photo IDs
	photoUUIDs := make([]uuid.UUID, 0, len(input.PhotoIDs))
	for _, id := range input.PhotoIDs {
		pid, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid photo ID: %s", id)
		}
		photoUUIDs = append(photoUUIDs, pid)
	}

	photos, err := u.photoRepo.FindByIDs(ctx, photoUUIDs)
	if err != nil {
		return nil, err
	}
	if len(photos) != len(input.PhotoIDs) {
		return nil, errors.New("one or more photos not found")
	}

	// ─ Gallery or Event context ─
	var pricePerPhoto float64
	var contextTitle string
	var photographerID uuid.UUID
	var galleryIDPtr *uuid.UUID
	var eventIDPtr *uuid.UUID

	if input.GalleryID != "" {
		galleryID, err := uuid.Parse(input.GalleryID)
		if err != nil {
			return nil, errors.New("invalid gallery_id")
		}
		gallery, err := u.galleryRepo.FindByID(ctx, galleryID)
		if err != nil {
			return nil, errors.New("gallery not found")
		}
		if !gallery.IsActive {
			return nil, errors.New("this gallery is no longer active")
		}
		if gallery.ExpiresAt != nil && gallery.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("this gallery has expired")
		}
		for _, photo := range photos {
			if photo.GalleryID == nil || *photo.GalleryID != galleryID {
				return nil, errors.New("photo does not belong to this gallery")
			}
		}
		pricePerPhoto = gallery.PricePerPhoto
		contextTitle = gallery.Title
		photographerID = gallery.PhotographerID
		galleryIDPtr = &galleryID

	} else {
		eventID, err := uuid.Parse(input.EventID)
		if err != nil {
			return nil, errors.New("invalid event_id")
		}
		event, err := u.eventRepo.FindByID(ctx, eventID)
		if err != nil {
			return nil, errors.New("event not found")
		}
		if !event.IsActive {
			return nil, errors.New("this event is no longer active")
		}
		pricePerPhoto = event.PricePerPhoto
		contextTitle = event.Title
		photographerID = event.PhotographerID
		eventIDPtr = &eventID
	}

	// ─ Commission calculation (business rule — done on backend) ─
	photographer, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return nil, err
	}

	totalAmount := pricePerPhoto * float64(len(photos))
	commissionRate := photographer.CommissionRate()
	platformFee := totalAmount * commissionRate
	netAmount := totalAmount - platformFee

	// ─ Create Order record ─
	order := &domain.Order{
		ID:            uuid.New(),
		BuyerEmail:    input.BuyerEmail,
		GalleryID:     galleryIDPtr,
		EventID:       eventIDPtr,
		TotalAmount:   totalAmount,
		PlatformFee:   platformFee,
		NetAmount:     netAmount,
		PaymentStatus: domain.PaymentPending,
	}

	for _, photo := range photos {
		order.OrderItems = append(order.OrderItems, domain.OrderItem{
			ID:      uuid.New(),
			PhotoID: photo.ID,
			Price:   pricePerPhoto,
		})
	}

	if err := u.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// ─ Generate QRIS via Midtrans Core API ─
	qrPayment, err := u.paymentService.GenerateQRCode(domain.QRCodePaymentRequest{
		OrderID:       order.ID.String(),
		Amount:        totalAmount,
		PaymentMethod: domain.PaymentMethodQRIS,
		CustomerEmail: input.BuyerEmail,
		Description:   fmt.Sprintf("Pembelian %d foto dari %s", len(photos), contextTitle),
	})
	if err != nil {
		// Don't fail order creation if QR fails — return order without QR
		// Client can poll /orders/:id to retry
		return &OrderSummary{Order: order}, fmt.Errorf("QR code generation failed: %w", err)
	}

	// Persist QR data to order record
	_ = u.orderRepo.UpdateQRCode(ctx, order.ID,
		qrPayment.QRCodeURL,
		qrPayment.QRCodeString,
		parseTime(qrPayment.ExpiresAt),
	)

	return &OrderSummary{Order: order, QRCodePayment: qrPayment}, nil
}

// ─── ProcessPaymentSuccess ─────────────────────────────────────────────────

// ProcessPaymentSuccess is called by the Midtrans webhook after settlement.
// Uses DB transaction to atomically update order + credit photographer wallet.
func (u *OrderUsecase) ProcessPaymentSuccess(ctx context.Context, gatewayOrderID, gatewayID, gatewayResp string) error {
	// Find order — Midtrans sends order_id which is our order UUID
	orderID, err := uuid.Parse(gatewayOrderID)
	if err != nil {
		return fmt.Errorf("invalid order ID from gateway: %s", gatewayOrderID)
	}
	order, err := u.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %s", gatewayOrderID)
	}

	if order.PaymentStatus == domain.PaymentPaid {
		return nil // Idempotent
	}

	// Resolve photographer and context title
	photographerID, contextTitle, err := u.resolveOrderContext(ctx, order)
	if err != nil {
		return err
	}

	photographer, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return err
	}

	return u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// 1. Update order payment status
		if err := u.orderRepo.UpdatePaymentStatus(txCtx, order.ID, domain.PaymentPaid, gatewayID, gatewayResp); err != nil {
			return err
		}

		// 2. Generate presigned download URLs (72h expiry)
		downloadExpiry := time.Now().Add(72 * time.Hour)
		var updatedItems []domain.OrderItem
		var mailItems []infra.DownloadItem

		for i, item := range order.OrderItems {
			photo, err := u.photoRepo.FindByID(txCtx, item.PhotoID)
			if err != nil {
				continue
			}

			s3Key := photo.S3KeyCompressed
			if s3Key == "" {
				s3Key = photo.S3Key
			}

			downloadURL, err := u.s3.GetPresignedDownloadURL(s3Key, 72*time.Hour)
			if err != nil {
				return fmt.Errorf("download URL for photo %s: %w", photo.ID, err)
			}

			order.OrderItems[i].DownloadURL = &downloadURL
			updatedItems = append(updatedItems, order.OrderItems[i])
			mailItems = append(mailItems, infra.DownloadItem{
				Filename:    photo.Filename,
				DownloadURL: downloadURL,
			})
		}

		if err := u.orderRepo.UpdateDownloadURLs(txCtx, order.ID, updatedItems); err != nil {
			return err
		}

		if err := u.orderRepo.UpdateDownloadExpiresAt(txCtx, order.ID, downloadExpiry); err != nil {
			return err
		}

		// 3. Credit photographer wallet (atomic)
		newBalance := photographer.WalletBalance + order.NetAmount
		walletTx := &domain.WalletTransaction{
			ID:             uuid.New(),
			PhotographerID: photographerID,
			OrderID:        &order.ID,
			Type:           domain.TransactionCredit,
			Amount:         order.NetAmount,
			Description:    fmt.Sprintf("Penjualan %d foto dari '%s'", len(order.OrderItems), contextTitle),
			BalanceAfter:   newBalance,
		}
		if err := u.walletRepo.AddTransaction(txCtx, walletTx); err != nil {
			return err
		}
		if err := u.userRepo.UpdateWalletBalance(txCtx, photographerID, order.NetAmount); err != nil {
			return err
		}

		// 4. Send download email (async — doesn't fail the DB transaction)
		// go func(data infra.DownloadLinkData) {
		// 	_ = u.mail.SendDownloadLinks(data)
		// }(infra.DownloadLinkData{
		// 	BuyerEmail:   order.BuyerEmail,
		// 	GalleryTitle: contextTitle,
		// 	OrderID:      order.ID.String(),
		// 	Items:        mailItems,
		// 	ExpiresAt:    downloadExpiry.Format("02 Jan 2006, 15:04 WIB"),
		// })

		return nil
	})
}

// GetOrder fetches a single order, verifying the buyer email for security.
func (u *OrderUsecase) GetOrder(ctx context.Context, orderID uuid.UUID, buyerEmail string) (*domain.Order, error) {
	order, err := u.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.BuyerEmail != buyerEmail {
		return nil, errors.New("order not found") // intentionally vague
	}
	return order, nil
}

// GetOrdersByEmail returns paginated orders for a buyer.
func (u *OrderUsecase) GetOrdersByEmail(ctx context.Context, email string, page, limit int) ([]domain.Order, int64, error) {
	return u.orderRepo.FindByBuyerEmail(ctx, email, page, limit)
}

// resolveOrderContext returns the photographerID and context title for an order.
func (u *OrderUsecase) resolveOrderContext(ctx context.Context, order *domain.Order) (uuid.UUID, string, error) {
	if order.GalleryID != nil {
		gallery, err := u.galleryRepo.FindByID(ctx, *order.GalleryID)
		if err != nil {
			return uuid.Nil, "", err
		}
		return gallery.PhotographerID, gallery.Title, nil
	}
	if order.EventID != nil {
		event, err := u.eventRepo.FindByID(ctx, *order.EventID)
		if err != nil {
			return uuid.Nil, "", err
		}
		return event.PhotographerID, event.Title, nil
	}
	return uuid.Nil, "", errors.New("order has no gallery or event")
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func parseTime(s string) *time.Time {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return &t
		}
	}
	return nil
}
