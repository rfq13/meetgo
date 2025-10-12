package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/webrtc-meeting/backend/internal/webrtc"
	"github.com/webrtc-meeting/backend/internal/websocket"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	// Parse command line flags
	var (
		port = flag.String("port", "8081", "Port for WebSocket server")
		host = flag.String("host", "0.0.0.0", "Host for WebSocket server")
		env  = flag.String("env", "development", "Environment (development/production)")
	)
	flag.Parse()

	// Setup logger
	logLevel := "info"
	if *env == "development" {
		logLevel = "debug"
	}

	loggerConfig := &logger.Config{
		Level:  logLevel,
		Format: "text",
		Output: "stdout",
	}

	logger.InitLogger(loggerConfig)
	logrus.Info("Starting WebSocket server...")

	// Setup Gin mode
	if *env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Janus client
	janusBaseURL := os.Getenv("JANUS_BASE_URL")
	if janusBaseURL == "" {
		janusBaseURL = "http://localhost:8088/janus"
	}

	janusAdminURL := os.Getenv("JANUS_ADMIN_URL")
	if janusAdminURL == "" {
		janusAdminURL = "http://localhost:8088/admin"
	}

	janusAPISecret := os.Getenv("JANUS_API_SECRET")
	janusAdminSecret := os.Getenv("JANUS_ADMIN_SECRET")

	janusClient := webrtc.NewJanusClient(janusBaseURL, janusAdminURL, janusAPISecret, janusAdminSecret)

	// Create WebSocket hub
	hub := websocket.NewHub()

	// Create signaling handler
	signalingHandler := webrtc.NewSignalingHandler(janusClient, hub)

	// Set signaling handler to hub
	hub.SignalingHandler = signalingHandler

	// Start hub in goroutine
	go hub.Run()

	// Create WebSocket handler
	wsHandler := websocket.NewHandler(hub)

	// Setup Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup WebSocket routes
	wsHandler.SetupRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"service":   "websocket-server",
			"version":   "1.0.0",
		})
	})

	// Create HTTP server
	serverAddr := fmt.Sprintf("%s:%s", *host, *port)
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in goroutine
	go func() {
		logrus.WithFields(logrus.Fields{
			"host": *host,
			"port": *port,
		}).Info("WebSocket server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down WebSocket server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}

	logrus.Info("WebSocket server stopped")
}
