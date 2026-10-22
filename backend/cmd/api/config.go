package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	appconfig "fotoin/backend/internal/config"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/handler"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/model"
	"fotoin/backend/internal/repository"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/validator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type application struct {
	server *echo.Echo
	close  func() error
}

func newApplication(ctx context.Context, cfg appconfig.Config) (_ *application, err error) {
	var closers []func() error
	defer func() {
		if err != nil {
			err = errors.Join(err, closeResources(closers))
		}
	}()

	db, err := infra.NewDatabase(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	closers = append(closers, closeDatabase(db))

	if err := infra.AutoMigrate(ctx, db); err != nil {
		return nil, fmt.Errorf("run database migrations: %w", err)
	}

	redisClient, err := infra.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("initialize redis: %w", err)
	}
	closers = append(closers, redisClient.Close)

	s3Client, err := infra.NewS3Client(ctx, cfg.S3)
	if err != nil {
		return nil, fmt.Errorf("initialize S3 client: %w", err)
	}

	txManager := infra.NewTxManager(db)
	userRepo := repository.NewUserRepository(db)
	galleryRepo := repository.NewGalleryRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	paymentService := infra.NewDuidkuService(cfg.Duitku)
	google := infra.NewGoogle(cfg.Google)
	mailClient := infra.NewMailClient(cfg.SMTP)

	authUsecase := usecase.NewAuthUsecase(userRepo, google)
	galleryUsecase := usecase.NewGalleryUsecase(galleryRepo, userRepo, s3Client)
	photoUsecase := usecase.NewPhotoUsecase(photoRepo, galleryRepo, eventRepo, userRepo, s3Client)
	orderUsecase := usecase.NewOrderUsecase(
		orderRepo, galleryRepo, eventRepo, photoRepo,
		userRepo, walletRepo, s3Client, mailClient, txManager, paymentService,
	)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, userRepo, txManager)
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, mailClient)
	subUsecase := usecase.NewSubscriptionUsecase(subRepo, userRepo, paymentService, txManager)
	adminUsecase := usecase.NewAdminUsecase(adminRepo, userRepo)

	authHandler := handler.NewAuthHandler(authUsecase)
	galleryHandler := handler.NewGalleryHandler(galleryUsecase)
	photoHandler := handler.NewPhotoHandler(photoUsecase, s3Client)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)
	eventHandler := handler.NewEventHandler(eventUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	subHandler := handler.NewSubscriptionHandler(subUsecase)

	// Start background retention cleaner (14 days basic retention policy)
	go startRetentionCleaner(db)

	server := newRouter(cfg, authHandler, galleryHandler, photoHandler, orderHandler, walletHandler, adminHandler, eventHandler, userHandler, subHandler)
	return &application{
		server: server,
		close:  func() error { return closeResources(closers) },
	}, nil
}

func newRouter(
	cfg appconfig.Config,
	authHandler *handler.AuthHandler,
	galleryHandler *handler.GalleryHandler,
	photoHandler *handler.PhotoHandler,
	orderHandler *handler.OrderHandler,
	walletHandler *handler.WalletHandler,
	adminHandler *handler.AdminHandler,
	eventHandler *handler.EventHandler,
	userHandler *handler.UserHandler,
	subHandler *handler.SubscriptionHandler,
) *echo.Echo {
	server := echo.New()
	server.Validator = validator.NewValidator()
	server.Use(echoMiddleware.Logger())
	server.Use(echoMiddleware.Recover())
	server.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	api := server.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/forgot-password", userHandler.ForgotPassword)
	auth.POST("/reset-password", userHandler.ResetPassword)
	auth.GET("/google", authHandler.GoogleLogin)
	auth.GET("/google/callback", authHandler.GoogleCallback)

	public := api.Group("/public")
	public.GET("/galleries", galleryHandler.ListPublic)
	public.GET("/galleries/:slug", galleryHandler.GetPublicGallery)
	public.GET("/galleries/:id/photos", photoHandler.ListPhotos)
	public.GET("/events", eventHandler.ListPublicEvents)
	public.GET("/events/:slug", eventHandler.GetEventBySlug)
	public.POST("/orders", orderHandler.CreateOrder)
	public.GET("/orders/:id", orderHandler.GetOrder)
	api.POST("/webhooks/payment", orderHandler.PaymentWebhook)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth)
	protected.GET("/me", adminHandler.GetMe)
	protected.PUT("/me/profile", userHandler.EditProfile)
	protected.PUT("/me/password", userHandler.UpdatePassword)
	protected.POST("/me/password", userHandler.CreatePassword)
	protected.DELETE("/me", userHandler.DeleteAccount)

	events := protected.Group("/events")
	events.GET("", eventHandler.ListMyEvents)
	events.POST("", eventHandler.CreateEvent)
	events.GET("/:id", eventHandler.GetEvent)
	events.PUT("/:id", eventHandler.UpdateEvent)
	events.DELETE("/:id", eventHandler.DeleteEvent)
	events.GET("/:id/photos", photoHandler.ListEventPhotos)
	events.POST("/:id/photos/upload-url", photoHandler.GetEventUploadURL)
	events.POST("/:id/photos/confirm", photoHandler.ConfirmEventUpload)
	events.POST("/:id/photos/proxy-upload", photoHandler.ProxyEventUpload)

	galleries := protected.Group("/galleries")
	galleries.GET("", galleryHandler.List)
	galleries.POST("", galleryHandler.Create)
	galleries.GET("/:id", galleryHandler.GetByID)
	galleries.PUT("/:id", galleryHandler.Update)
	galleries.DELETE("/:id", galleryHandler.Delete)
	galleries.GET("/:id/photos", photoHandler.ListPhotos)
	galleries.POST("/:id/photos/upload-url", photoHandler.GetUploadURL)
	galleries.POST("/:id/photos/confirm", photoHandler.ConfirmUpload)
	galleries.POST("/:id/photos/proxy-upload", photoHandler.ProxyUpload)

	protected.DELETE("/photos/:id", photoHandler.DeletePhoto)
	protected.GET("/wallet", walletHandler.GetWallet)
	protected.POST("/wallet/withdrawal", walletHandler.RequestWithdrawal)

	// Subscription endpoints for photographers
	subscriptions := protected.Group("/subscriptions")
	subscriptions.GET("/me", subHandler.GetMySubscription)
	subscriptions.POST("/upgrade", subHandler.Upgrade)
	subscriptions.POST("/:id/simulate", subHandler.SimulatePayment)

	clientRoutes := protected.Group("/client")
	clientRoutes.Use(middleware.RequireRole(domain.RoleClient, domain.RoleSuperAdmin))
	clientRoutes.GET("/orders", orderHandler.GetMyOrders)

	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(domain.RoleSuperAdmin))
	admin.GET("/users", adminHandler.ListUsers)
	admin.PATCH("/users/:id/suspend", adminHandler.SuspendUser)
	admin.GET("/revenue", adminHandler.GetRevenue)
	admin.GET("/withdrawals", walletHandler.ListWithdrawals)
	admin.PATCH("/withdrawals/:id", walletHandler.ApproveWithdrawal)
	admin.GET("/subscriptions", subHandler.ListAdminSubscriptions)
	admin.POST("/subscriptions/grant", subHandler.AdminGrantPremium)

	return server
}

func startRetentionCleaner(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	clean := func() {
		now := time.Now()
		var count int64
		err := db.Model(&model.GalleryModel{}).
			Where("is_active = true AND expires_at IS NOT NULL AND expires_at <= ?", now).
			Count(&count).Error
		if err == nil && count > 0 {
			db.Model(&model.GalleryModel{}).
				Where("is_active = true AND expires_at IS NOT NULL AND expires_at <= ?", now).
				Update("is_active", false)
			log.Printf("[RetentionCleaner] Deactivated %d expired galleries according to 14-day retention rule", count)
		}
	}
	// Initial check on startup
	clean()
	for range ticker.C {
		clean()
	}
}

func closeDatabase(db *gorm.DB) func() error {
	return func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("get database connection for closing: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("close database connection: %w", err)
		}
		return nil
	}
}

func closeResources(closers []func() error) error {
	var closeErrors []error
	for i := len(closers) - 1; i >= 0; i-- {
		if err := closers[i](); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return errors.Join(closeErrors...)
}
