package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/webrtc-meeting/backend/internal/api"
	"github.com/webrtc-meeting/backend/internal/auth"
	"github.com/webrtc-meeting/backend/internal/config"
	"github.com/webrtc-meeting/backend/internal/database"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(&logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: "stdout",
	})

	log.Info("Starting WebRTC Meeting backend server")

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize database")
	}

	// Initialize services and router
	authService := auth.NewService(db.DB, cfg, log)
	router, err := api.InitializeRouter(db.DB, log, authService)
	if err != nil {
		log.WithError(err).Fatal("Failed to initialize router")
	}

	engine := router.SetupRouter()

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("HTTP server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	} else {
		log.Info("Server shutdown completed")
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.WithError(err).Error("Failed to close database connection")
	}

	log.Info("Bye")
}
