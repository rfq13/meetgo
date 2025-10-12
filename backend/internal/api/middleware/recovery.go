package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// RecoveryConfig konfigurasi recovery middleware
type RecoveryConfig struct {
	// Logger untuk logging
	Logger *logger.Logger

	// Stack trace configuration
	StackSize    int
	DisableStack bool

	// Response configuration
	AllowerPanic bool

	// Custom recovery function
	CustomRecoveryFunc gin.HandlerFunc

	// Request body configuration
	MaxRequestBodySize int64
}

// DefaultRecoveryConfig konfigurasi default
func DefaultRecoveryConfig(log *logger.Logger) RecoveryConfig {
	return RecoveryConfig{
		Logger:             log,
		StackSize:          4 * 1024, // 4KB
		DisableStack:       false,
		AllowerPanic:       false,
		MaxRequestBodySize: 1024, // 1KB
	}
}

// PanicInfo informasi panic yang terjadi
type PanicInfo struct {
	Message   string                 `json:"message"`
	Stack     string                 `json:"stack,omitempty"`
	Request   map[string]interface{} `json:"request,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Recovery middleware untuk recovery dari panic
func Recovery(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic
				config.logPanic(c, err)

				// Execute custom recovery function if provided
				if config.CustomRecoveryFunc != nil {
					config.CustomRecoveryFunc(c)
					return
				}

				// Default recovery response
				if !c.Writer.Written() {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal server error",
						"message": "Something went wrong, please try again later",
						"code":    "INTERNAL_ERROR",
					})
				}
			}
		}()

		c.Next()
	}
}

// logPanic logging informasi panic
func (config RecoveryConfig) logPanic(c *gin.Context, err interface{}) {
	stack := ""
	if !config.DisableStack {
		stack = string(debug.Stack())
		if config.StackSize > 0 && len(stack) > config.StackSize {
			stack = stack[:config.StackSize]
		}
	}

	// Get request information
	requestInfo := config.getRequestInfo(c)

	panicInfo := PanicInfo{
		Message:   fmt.Sprintf("%v", err),
		Stack:     stack,
		Request:   requestInfo,
		Timestamp: time.Now(),
	}

	// Log panic
	if config.Logger != nil {
		config.Logger.WithFields(map[string]interface{}{
			"panic_message": panicInfo.Message,
			"request_info":  panicInfo.Request,
			"stack_trace":   panicInfo.Stack,
		}).Error("Panic recovered")
	}
}

// getRequestInfo mengambil informasi request
func (config RecoveryConfig) getRequestInfo(c *gin.Context) map[string]interface{} {
	requestInfo := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.RawQuery,
		"ip":         c.ClientIP(),
		"user_agent": c.Request.UserAgent(),
		"headers":    config.getSafeHeaders(c),
	}

	// Add request body if it's small enough
	if c.Request.Body != nil && config.MaxRequestBodySize > 0 {
		body := config.getRequestBody(c)
		if body != "" {
			requestInfo["body"] = body
		}
	}

	// Add user info if available
	if userID, exists := c.Get("user_id"); exists {
		requestInfo["user_id"] = userID
	}

	return requestInfo
}

// getSafeHeaders mengambil headers yang aman untuk di-log
func (config RecoveryConfig) getSafeHeaders(c *gin.Context) map[string]string {
	safeHeaders := map[string]string{
		"Content-Type":    c.GetHeader("Content-Type"),
		"Accept":          c.GetHeader("Accept"),
		"User-Agent":      c.GetHeader("User-Agent"),
		"Referer":         c.GetHeader("Referer"),
		"X-Forwarded-For": c.GetHeader("X-Forwarded-For"),
		"X-Real-IP":       c.GetHeader("X-Real-IP"),
	}

	// Add Authorization header without token
	if auth := c.GetHeader("Authorization"); auth != "" {
		if parts := strings.SplitN(auth, " ", 2); len(parts) == 2 {
			safeHeaders["Authorization"] = parts[0] + " [REDACTED]"
		} else {
			safeHeaders["Authorization"] = "[REDACTED]"
		}
	}

	return safeHeaders
}

// getRequestBody mengambil request body dengan batasan ukuran
func (config RecoveryConfig) getRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// Read body
	bodyBytes, err := io.ReadAll(io.LimitReader(c.Request.Body, config.MaxRequestBodySize))
	if err != nil {
		return ""
	}

	// Restore body for next handler
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Try to parse as JSON for better logging
	var jsonData interface{}
	if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
		if formatted, err := json.Marshal(jsonData); err == nil {
			return string(formatted)
		}
	}

	return string(bodyBytes)
}

// CustomRecovery membuat custom recovery function
func CustomRecovery(recoveryFunc func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				recoveryFunc(c, err)
			}
		}()
		c.Next()
	}
}

// SimpleRecovery middleware recovery sederhana
func SimpleRecovery(log *logger.Logger) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.DisableStack = true
	return Recovery(config)
}

// DetailedRecovery middleware recovery dengan detail stack trace
func DetailedRecovery(log *logger.Logger) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.StackSize = 8 * 1024 // 8KB
	return Recovery(config)
}

// ProductionRecovery middleware recovery untuk production
func ProductionRecovery(log *logger.Logger) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.DisableStack = true
	config.MaxRequestBodySize = 0 // Don't log request body in production
	return Recovery(config)
}

// DevelopmentRecovery middleware recovery untuk development
func DevelopmentRecovery(log *logger.Logger) gin.HandlerFunc {
	config := DefaultRecoveryConfig(log)
	config.StackSize = 16 * 1024     // 16KB
	config.MaxRequestBodySize = 2048 // 2KB
	return Recovery(config)
}

// ErrorHandler middleware untuk error handling tambahan
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors that occurred during request processing
		for _, err := range c.Errors {
			log.WithFields(map[string]interface{}{
				"error":  err.Error(),
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"ip":     c.ClientIP(),
			}).Error("Request error")
		}

		// Handle specific error types
		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last()

			// Don't override response if already written
			if !c.Writer.Written() {
				switch {
				case strings.Contains(lastErr.Error(), "bind"):
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   "Bad request",
						"message": "Invalid request format",
						"code":    "INVALID_REQUEST",
					})
				case strings.Contains(lastErr.Error(), "validation"):
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"error":   "Validation error",
						"message": "Request validation failed",
						"code":    "VALIDATION_ERROR",
					})
				case strings.Contains(lastErr.Error(), "unauthorized"):
					c.JSON(http.StatusUnauthorized, gin.H{
						"error":   "Unauthorized",
						"message": "Authentication required",
						"code":    "UNAUTHORIZED",
					})
				case strings.Contains(lastErr.Error(), "forbidden"):
					c.JSON(http.StatusForbidden, gin.H{
						"error":   "Forbidden",
						"message": "Access denied",
						"code":    "FORBIDDEN",
					})
				case strings.Contains(lastErr.Error(), "not found"):
					c.JSON(http.StatusNotFound, gin.H{
						"error":   "Not found",
						"message": "Resource not found",
						"code":    "NOT_FOUND",
					})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Internal server error",
						"message": "Something went wrong",
						"code":    "INTERNAL_ERROR",
					})
				}
			}
		}
	}
}
