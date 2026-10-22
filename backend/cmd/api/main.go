package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fotoin/backend/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := config.LoadConfig()
	app, err := newApplication(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}
	defer func() {
		if err := app.close(); err != nil {
			log.Printf("Failed to close application resources: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("FotoIN API running on :%s", cfg.Port)
		serverErr <- app.server.Start(":" + cfg.Port)
	}()

	serverStopped := false
	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErr:
		serverStopped = true
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("API server failed: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := app.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}

	if !serverStopped {
		// Give the server goroutine time to report its final error before cleanup.
		select {
		case err := <-serverErr:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("API server stopped with error: %v", err)
			}
		case <-time.After(cfg.ShutdownTimeout):
			log.Println("Timed out waiting for API server to stop")
		}
	}
}
