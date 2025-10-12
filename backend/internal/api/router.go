package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/webrtc-meeting/backend/internal/api/middleware"
	"github.com/webrtc-meeting/backend/internal/auth"
	"github.com/webrtc-meeting/backend/internal/room"
	"github.com/webrtc-meeting/backend/internal/user"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Router struct untuk menyimpan semua dependencies
type Router struct {
	db          *gorm.DB
	logger      *logger.Logger
	authHandler *auth.Handler
	userHandler *user.Handler
	roomHandler *room.Handler
}

// NewRouter membuat router baru dengan semua dependencies
func NewRouter(
	db *gorm.DB,
	log *logger.Logger,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	roomHandler *room.Handler,
) *Router {
	return &Router{
		db:          db,
		logger:      log,
		authHandler: authHandler,
		userHandler: userHandler,
		roomHandler: roomHandler,
	}
}

// SetupRouter mengkonfigurasi semua routes dan middleware
func (r *Router) SetupRouter() *gin.Engine {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add global middleware
	r.setupGlobalMiddleware(router)

	// Setup API routes
	r.setupAPIRoutes(router)

	// Setup health check
	r.setupHealthCheck(router)

	// Setup static files (if needed)
	r.setupStaticFiles(router)

	return router
}

// setupGlobalMiddleware mengatur middleware global
func (r *Router) setupGlobalMiddleware(router *gin.Engine) {
	// Recovery middleware
	router.Use(middleware.DetailedRecovery(r.logger))

	// Logger middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		r.logger.LogRequest(
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
		return ""
	}))

	// Error handler middleware
	router.Use(middleware.ErrorHandler(r.logger))

	// CORS middleware
	corsConfig := middleware.DefaultCORSConfig()
	router.Use(middleware.CORS(corsConfig))

	// Rate limiting middleware
	router.Use(middleware.APIRateLimit(r.logger))
}

// setupAPIRoutes mengatur API routes
func (r *Router) setupAPIRoutes(router *gin.Engine) {
	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (with stricter rate limiting)
		auth := v1.Group("/auth")
		auth.Use(middleware.AuthRateLimit(r.logger))
		{
			r.authHandler.RegisterRoutes(auth)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(r.authHandler.AuthMiddleware())
		{
			// User routes
			r.userHandler.RegisterRoutes(protected)

			// Room routes
			r.roomHandler.RegisterRoutes(protected)
		}

		// Admin routes (require admin role)
		admin := v1.Group("/admin")
		admin.Use(r.authHandler.AuthMiddleware())
		admin.Use(r.authHandler.AdminMiddleware())
		{
			r.setupAdminRoutes(admin)
		}

		// WebRTC routes
		webrtc := v1.Group("/webrtc")
		webrtc.Use(r.authHandler.AuthMiddleware())
		{
			r.setupWebRTCRoutes(webrtc)
		}

		// Public routes (optional authentication)
		public := v1.Group("/public")
		{
			r.setupPublicRoutes(public)
		}
	}
}

// setupAdminRoutes mengatur admin routes
func (r *Router) setupAdminRoutes(admin *gin.RouterGroup) {
	// Get users (admin only)
	admin.GET("/users", r.adminGetUsers)

	// Update user status (admin only)
	admin.PUT("/users/:userId/status", r.adminUpdateUserStatus)

	// Get rooms (admin only)
	admin.GET("/rooms", r.adminGetRooms)

	// Get system statistics (admin only)
	admin.GET("/statistics", r.adminGetStatistics)

	// System health check (admin only)
	admin.GET("/health", r.adminHealthCheck)
}

// setupWebRTCRoutes mengatur WebRTC routes
func (r *Router) setupWebRTCRoutes(webrtc *gin.RouterGroup) {
	// Get ICE servers
	webrtc.GET("/ice-servers", r.getICEServers)

	// Recording routes
	recording := webrtc.Group("/recording")
	{
		recording.GET("", r.getRecordings)
		recording.GET("/:recordingId", r.getRecording)
	}

	// Room recording controls
	rooms := webrtc.Group("/rooms")
	{
		rooms.POST("/:roomId/recording/start", r.startRecording)
		rooms.POST("/:roomId/recording/stop", r.stopRecording)
	}
}

// setupPublicRoutes mengatur public routes
func (r *Router) setupPublicRoutes(public *gin.RouterGroup) {
	// Public room listing (with optional auth)
	public.GET("/rooms", r.roomHandler.GetRooms)

	// Public room details (with optional auth)
	public.GET("/rooms/:roomId", r.roomHandler.GetRoom)

	// System info
	public.GET("/info", r.getSystemInfo)

	// API documentation
	public.GET("/docs", r.getAPIDocumentation)
}

// setupHealthCheck mengatur health check routes
func (r *Router) setupHealthCheck(router *gin.Engine) {
	router.GET("/health", r.healthCheck)
	router.GET("/ping", r.ping)
	router.GET("/ready", r.readinessCheck)
	router.GET("/live", r.livenessCheck)
}

// setupStaticFiles mengatur static file serving
func (r *Router) setupStaticFiles(router *gin.Engine) {
	// Serve static files if needed
	// router.Static("/static", "./static")
	// router.StaticFS("/uploads", gin.Dir("./uploads"))
}

// Health Check Handlers

func (r *Router) healthCheck(c *gin.Context) {
	status := r.checkSystemHealth()

	if status["status"] == "healthy" {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"checks":    status,
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "unhealthy",
			"timestamp": time.Now().UTC(),
			"checks":    status,
		})
	}
}

func (r *Router) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "pong",
		"timestamp": time.Now().UTC(),
	})
}

func (r *Router) readinessCheck(c *gin.Context) {
	// Check if application is ready to serve traffic
	ready := r.checkReadiness()

	if ready {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now().UTC(),
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"timestamp": time.Now().UTC(),
		})
	}
}

func (r *Router) livenessCheck(c *gin.Context) {
	// Check if application is alive
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}

// Admin Handlers (placeholders - to be implemented)

func (r *Router) adminGetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin get users - to be implemented",
	})
}

func (r *Router) adminUpdateUserStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin update user status - to be implemented",
	})
}

func (r *Router) adminGetRooms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin get rooms - to be implemented",
	})
}

func (r *Router) adminGetStatistics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin get statistics - to be implemented",
	})
}

func (r *Router) adminHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin health check - to be implemented",
	})
}

// WebRTC Handlers (placeholders - to be implemented)

func (r *Router) getICEServers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get ICE servers - to be implemented",
	})
}

func (r *Router) getRecordings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get recordings - to be implemented",
	})
}

func (r *Router) getRecording(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get recording - to be implemented",
	})
}

func (r *Router) startRecording(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Start recording - to be implemented",
	})
}

func (r *Router) stopRecording(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Stop recording - to be implemented",
	})
}

// Public Handlers (placeholders - to be implemented)

func (r *Router) getSystemInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get system info - to be implemented",
	})
}

func (r *Router) getAPIDocumentation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Get API documentation - to be implemented",
	})
}

// Health Check Helper Functions

func (r *Router) checkSystemHealth() map[string]interface{} {
	checks := make(map[string]interface{})
	overallStatus := "healthy"

	// Check database
	dbStatus := r.checkDatabaseHealth()
	checks["database"] = dbStatus
	if dbStatus["status"] != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check other services here
	// Example: Redis, external APIs, etc.

	checks["status"] = overallStatus
	return checks
}

func (r *Router) checkDatabaseHealth() map[string]interface{} {
	// Get database stats
	sqlDB, err := r.db.DB()
	if err != nil {
		return map[string]interface{}{
			"status":  "unhealthy",
			"message": "Failed to get database connection",
			"error":   err.Error(),
		}
	}

	// Ping database
	if err := sqlDB.Ping(); err != nil {
		return map[string]interface{}{
			"status":  "unhealthy",
			"message": "Database ping failed",
			"error":   err.Error(),
		}
	}

	// Get connection stats
	stats := sqlDB.Stats()

	return map[string]interface{}{
		"status":               "healthy",
		"message":              "Database connection OK",
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

func (r *Router) checkReadiness() bool {
	// Check if all critical dependencies are ready
	if err := r.checkDatabaseReadiness(); err != nil {
		r.logger.WithError(err).Error("Database not ready")
		return false
	}

	// Check other readiness criteria here
	// Example: external services, configuration, etc.

	return true
}

func (r *Router) checkDatabaseReadiness() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// Helper function to inject auth middleware into room and user handlers
func (r *Router) injectAuthMiddleware() {
	// Set auth middleware in room handler
	r.roomHandler.SetAuthMiddleware(r.authHandler.AuthMiddleware())
	r.roomHandler.SetOptionalAuthMiddleware(r.authHandler.OptionalAuthMiddleware())

	// User handler already uses the auth handler directly
}

// Initialize router dengan dependencies yang lengkap
func InitializeRouter(
	db *gorm.DB,
	log *logger.Logger,
	authService *auth.Service,
) (*Router, error) {
	// Create handlers
	authHandler := auth.NewHandler(authService, log)
	userService := user.NewService(db, log)
	userHandler := user.NewHandler(userService, log)
	roomService := room.NewService(db, log)
	roomHandler := room.NewHandler(roomService, log)

	// Create router
	router := NewRouter(db, log, authHandler, userHandler, roomHandler)

	// Inject auth middleware
	router.injectAuthMiddleware()

	return router, nil
}
