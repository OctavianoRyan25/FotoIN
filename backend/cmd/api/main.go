package main

import (
	"log"
	"os"
	"strings"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/handler"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/repository"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/validator"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Infrastructure
	db := infra.NewDatabase()
	infra.AutoMigrate(db)

	// Init Transaction Manager
	txManager := infra.NewTxManager(db)

	redis := infra.NewRedisClient()
	s3Client := infra.NewS3Client()
	mailClient := infra.NewMailClient()
	google := infra.NewGoogle()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	galleryRepo := repository.NewGalleryRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Payment service — Duitku QRIS integration
	paymentService := infra.NewDuidkuService()

	// Usecases
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

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	galleryHandler := handler.NewGalleryHandler(galleryUsecase)
	photoHandler := handler.NewPhotoHandler(photoUsecase, s3Client)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	walletHandler := handler.NewWalletHandler(walletUsecase)
	adminHandler := handler.NewAdminHandler(userRepo, db, redis.Client())
	eventHandler := handler.NewEventHandler(eventUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	// Echo setup
	e := echo.New()
	e.Validator = validator.NewValidator()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")
	if len(origins) == 0 || origins[0] == "" {
		origins = []string{"*"}
	}

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	// ─── Routes ───────────────────────────────────────────────────────────────
	api := e.Group("/api/v1")

	// ── Auth (public) ──
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/forgot-password", userHandler.ForgotPassword)
	auth.POST("/reset-password", userHandler.ResetPassword)
	// Google OAuth — browser navigates here
	auth.GET("/google", authHandler.GoogleLogin)
	auth.GET("/google/callback", authHandler.GoogleCallback)

	// ── Public (no auth) ──
	pub := api.Group("/public")
	pub.GET("/galleries", galleryHandler.ListPublic)
	pub.GET("/galleries/:slug", galleryHandler.GetPublicGallery)
	pub.GET("/galleries/:id/photos", photoHandler.ListPhotos)
	pub.GET("/events", eventHandler.ListPublicEvents)
	pub.GET("/events/:slug", eventHandler.GetEventBySlug)
	pub.POST("/orders", orderHandler.CreateOrder)
	pub.GET("/orders/:id", orderHandler.GetOrder)

	// ── Payment webhook (Duidku notifies here) ──
	api.POST("/webhooks/payment", orderHandler.PaymentWebhook)

	// ── Protected (JWT required) ──
	protected := api.Group("")
	protected.Use(middleware.JWTAuth)

	protected.GET("/me", adminHandler.GetMe)

	// User profile & account management (all authenticated roles)
	protected.PUT("/me/profile", userHandler.EditProfile)
	protected.PUT("/me/password", userHandler.UpdatePassword)
	protected.POST("/me/password", userHandler.CreatePassword)
	protected.DELETE("/me", userHandler.DeleteAccount)

	// Event management (photographer)
	events := protected.Group("/events")
	events.GET("", eventHandler.ListMyEvents)
	events.POST("", eventHandler.CreateEvent)
	events.GET("/:id", eventHandler.GetEvent)
	events.PUT("/:id", eventHandler.UpdateEvent)
	events.DELETE("/:id", eventHandler.DeleteEvent)
	// Photos under event
	events.GET("/:id/photos", photoHandler.ListEventPhotos)
	events.POST("/:id/photos/upload-url", photoHandler.GetEventUploadURL)
	events.POST("/:id/photos/confirm", photoHandler.ConfirmEventUpload)
	// Proxy upload: browser → backend → S3 (no S3 CORS needed)
	events.POST("/:id/photos/proxy-upload", photoHandler.ProxyEventUpload)

	// Gallery management (photographer)
	galleries := protected.Group("/galleries")
	galleries.GET("", galleryHandler.List)
	galleries.POST("", galleryHandler.Create)
	galleries.GET("/:id", galleryHandler.GetByID)
	galleries.PUT("/:id", galleryHandler.Update)
	galleries.DELETE("/:id", galleryHandler.Delete)
	galleries.GET("/:id/photos", photoHandler.ListPhotos)
	galleries.POST("/:id/photos/upload-url", photoHandler.GetUploadURL)
	galleries.POST("/:id/photos/confirm", photoHandler.ConfirmUpload)
	// Proxy upload: browser → backend → S3 (no S3 CORS needed)
	galleries.POST("/:id/photos/proxy-upload", photoHandler.ProxyUpload)

	// Photo management
	photos := protected.Group("/photos")
	photos.DELETE("/:id", photoHandler.DeletePhoto)

	// Wallet
	wallet := protected.Group("/wallet")
	wallet.GET("", walletHandler.GetWallet)
	wallet.POST("/withdrawal", walletHandler.RequestWithdrawal)

	// Client — authenticated client routes
	clientRoutes := protected.Group("/client")
	clientRoutes.Use(middleware.RequireRole(domain.RoleClient, domain.RoleSuperAdmin))
	clientRoutes.GET("/orders", orderHandler.GetMyOrders)

	// Super Admin
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole(domain.RoleSuperAdmin))
	admin.GET("/users", adminHandler.ListUsers)
	admin.PATCH("/users/:id/suspend", adminHandler.SuspendUser)
	admin.GET("/revenue", adminHandler.GetRevenue)
	admin.GET("/withdrawals", walletHandler.ListWithdrawals)
	admin.PATCH("/withdrawals/:id", walletHandler.ApproveWithdrawal)

	// Start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 FotoIN API running on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
