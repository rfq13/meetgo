package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/webrtc-meeting/backend/models"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Handler struct untuk authentication handler
type Handler struct {
	service *Service
	logger  *logger.Logger
}

// NewHandler membuat authentication handler baru
func NewHandler(service *Service, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  log,
	}
}

// RegisterRoutes registrasi routes untuk authentication
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		auth.POST("/refresh", h.RefreshToken)
		auth.POST("/logout", h.AuthMiddleware(), h.Logout)
		auth.GET("/profile", h.AuthMiddleware(), h.GetProfile)
		auth.PUT("/password", h.AuthMiddleware(), h.ChangePassword)
	}
}

// Login handler untuk login endpoint
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	clientIP := c.ClientIP()
	response, err := h.service.Login(&req, clientIP)
	if err != nil {
		h.logger.WithError(err).WithField("email", req.Email).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithUserID(response.User.ID.String()).Info("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// Register handler untuk register endpoint
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid register request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	clientIP := c.ClientIP()
	user, err := h.service.Register(&req, clientIP)
	if err != nil {
		h.logger.WithError(err).WithField("email", req.Email).Error("Registration failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithUserID(user.ID.String()).Info("User registered successfully")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    user,
	})
}

// RefreshToken handler untuk refresh token endpoint
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid refresh token request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	clientIP := c.ClientIP()
	response, err := h.service.RefreshToken(&req, clientIP)
	if err != nil {
		h.logger.WithError(err).Error("Token refresh failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithUserID(response.User.ID.String()).Info("Token refreshed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

// Logout handler untuk logout endpoint
func (h *Handler) Logout(c *gin.Context) {
	token := h.extractTokenFromHeader(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No token provided",
		})
		return
	}

	clientIP := c.ClientIP()
	if err := h.service.Logout(token, clientIP); err != nil {
		h.logger.WithError(err).Error("Logout failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
		})
		return
	}

	h.logger.Info("User logged out successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// GetProfile handler untuk get profile endpoint
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	profile, err := h.service.GetProfile(userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get profile")
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Profile retrieved successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    profile,
	})
}

// ChangePassword handler untuk change password endpoint
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid change password request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	if err := h.service.ChangePassword(userUUID, &req); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Password change failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Password changed successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// AuthMiddleware middleware untuk authentication
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := h.extractTokenFromHeader(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			c.Abort()
			return
		}

		user, err := h.service.ValidateToken(token)
		if err != nil {
			h.logger.WithError(err).WithField("token", token).Warn("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", user.ID.String())
		c.Set("user_email", user.Email)
		c.Set("user_role", string(user.Role))
		c.Set("user", user)

		c.Next()
	}
}

// AdminMiddleware middleware untuk admin only
func (h *Handler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		if userRole != string(models.UserRoleAdmin) {
			h.logger.WithField("user_role", userRole).Warn("Unauthorized admin access attempt")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ModeratorMiddleware middleware untuk moderator dan admin
func (h *Handler) ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		if userRole != string(models.UserRoleAdmin) && userRole != string(models.UserRoleModerator) {
			h.logger.WithField("user_role", userRole).Warn("Unauthorized moderator access attempt")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Moderator access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddleware middleware untuk optional authentication
func (h *Handler) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := h.extractTokenFromHeader(c.GetHeader("Authorization"))
		if token != "" {
			user, err := h.service.ValidateToken(token)
			if err == nil {
				// Set user context if token is valid
				c.Set("user_id", user.ID.String())
				c.Set("user_email", user.Email)
				c.Set("user_role", string(user.Role))
				c.Set("user", user)
			}
		}

		c.Next()
	}
}

// extractTokenFromHeader extracts token from Authorization header
func (h *Handler) extractTokenFromHeader(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetCurrentUser helper function to get current user from context
func (h *Handler) GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	currentUser, ok := user.(*models.User)
	if !ok {
		return nil, false
	}

	return currentUser, true
}

// GetCurrentUser helper function to get current user ID from context
func (h *Handler) GetCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.Nil, false
	}

	return userUUID, true
}

// PaginationParams struct for pagination
type PaginationParams struct {
	Page    int `form:"page" binding:"min=1"`
	PerPage int `form:"per_page" binding:"min=1,max=100"`
}

// GetPaginationParams helper function to get pagination params
func (h *Handler) GetPaginationParams(c *gin.Context) PaginationParams {
	params := PaginationParams{
		Page:    1,
		PerPage: 20,
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 && perPage <= 100 {
			params.PerPage = perPage
		}
	}

	return params
}

// SuccessResponse helper function for success response
func (h *Handler) SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
}

// ErrorResponse helper function for error response
func (h *Handler) ErrorResponse(c *gin.Context, statusCode int, message string, details interface{}) {
	response := gin.H{
		"error": message,
	}

	if details != nil {
		response["details"] = details
	}

	c.JSON(statusCode, response)
}

// PaginatedResponse helper function for paginated response
func (h *Handler) PaginatedResponse(c *gin.Context, message string, data interface{}, total int64, params PaginationParams) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
		"pagination": gin.H{
			"page":        params.Page,
			"per_page":    params.PerPage,
			"total":       total,
			"total_pages": (total + int64(params.PerPage) - 1) / int64(params.PerPage),
		},
	})
}
