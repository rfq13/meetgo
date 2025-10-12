package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Handler struct untuk user handler
type Handler struct {
	service *Service
	logger  *logger.Logger
}

// NewHandler membuat user handler baru
func NewHandler(service *Service, log *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  log,
	}
}

// RegisterRoutes registrasi routes untuk user management
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		// Profile endpoints
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.PUT("/password", h.ChangePassword)

		// Contact endpoints
		users.GET("/contacts", h.GetContacts)
		users.POST("/contacts", h.AddContact)
		users.DELETE("/contacts/:contactId", h.RemoveContact)
		users.GET("/search", h.SearchUsers)

		// Settings endpoints
		users.GET("/settings", h.GetSettings)
		users.PUT("/settings", h.UpdateSettings)

		// Notification endpoints
		users.GET("/notifications", h.GetNotifications)
		users.PUT("/notifications/:notificationId/read", h.MarkNotificationAsRead)
		users.PUT("/notifications/read-all", h.MarkAllNotificationsAsRead)
		users.DELETE("/notifications/:notificationId", h.DeleteNotification)

		// Stats endpoint
		users.GET("/stats", h.GetUserStats)

		// Account management
		users.PUT("/deactivate", h.DeactivateAccount)
	}
}

// GetProfile handler untuk get profile endpoint
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	profile, err := h.service.GetProfile(userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get profile")
		h.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Profile retrieved successfully")
	h.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// UpdateProfile handler untuk update profile endpoint
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update profile request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	profile, err := h.service.UpdateProfile(userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to update profile")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Profile updated successfully")
	h.SuccessResponse(c, "Profile updated successfully", profile)
}

// ChangePassword handler untuk change password endpoint
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid change password request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.ChangePassword(userUUID, &req); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Password change failed")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Password changed successfully")
	h.SuccessResponse(c, "Password changed successfully", nil)
}

// GetContacts handler untuk get contacts endpoint
func (h *Handler) GetContacts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	params := h.GetPaginationParams(c)
	search := c.Query("search")

	contacts, total, err := h.service.GetContacts(userUUID, params.Page, params.PerPage, search)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get contacts")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Contacts retrieved successfully")
	h.PaginatedResponse(c, "Contacts retrieved successfully", contacts, total, params)
}

// AddContact handler untuk add contact endpoint
func (h *Handler) AddContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req AddContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid add contact request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	contactUUID, err := uuid.Parse(req.ContactID)
	if err != nil {
		h.logger.WithError(err).Error("Invalid contact ID")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", nil)
		return
	}

	contact, err := h.service.AddContact(userUUID, contactUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to add contact")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("contact_id", contactUUID.String()).Info("Contact added successfully")
	h.SuccessResponse(c, "Contact added successfully", contact)
}

// RemoveContact handler untuk remove contact endpoint
func (h *Handler) RemoveContact(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	contactUUID, err := uuid.Parse(c.Param("contactId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid contact ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid contact ID", nil)
		return
	}

	if err := h.service.RemoveContact(userUUID, contactUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to remove contact")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("contact_id", contactUUID.String()).Info("Contact removed successfully")
	h.SuccessResponse(c, "Contact removed successfully", nil)
}

// SearchUsers handler untuk search users endpoint
func (h *Handler) SearchUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	query := c.Query("q")
	if query == "" {
		h.ErrorResponse(c, http.StatusBadRequest, "Search query is required", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	params := h.GetPaginationParams(c)

	users, total, err := h.service.SearchUsers(userUUID, query, params.Page, params.PerPage)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to search users")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("query", query).Info("Users searched successfully")
	h.PaginatedResponse(c, "Users retrieved successfully", users, total, params)
}

// GetSettings handler untuk get settings endpoint
func (h *Handler) GetSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	settings, err := h.service.GetSettings(userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get settings")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Settings retrieved successfully")
	h.SuccessResponse(c, "Settings retrieved successfully", settings)
}

// UpdateSettings handler untuk update settings endpoint
func (h *Handler) UpdateSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update settings request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	settings, err := h.service.UpdateSettings(userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to update settings")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Settings updated successfully")
	h.SuccessResponse(c, "Settings updated successfully", settings)
}

// GetNotifications handler untuk get notifications endpoint
func (h *Handler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	params := h.GetPaginationParams(c)

	// Parse is_read parameter
	var isRead *bool
	if isReadStr := c.Query("is_read"); isReadStr != "" {
		if isReadBool, err := strconv.ParseBool(isReadStr); err == nil {
			isRead = &isReadBool
		}
	}

	notifications, total, err := h.service.GetNotifications(userUUID, params.Page, params.PerPage, isRead)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get notifications")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("Notifications retrieved successfully")
	h.PaginatedResponse(c, "Notifications retrieved successfully", notifications, total, params)
}

// MarkNotificationAsRead handler untuk mark notification as read endpoint
func (h *Handler) MarkNotificationAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	notificationUUID, err := uuid.Parse(c.Param("notificationId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid notification ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", nil)
		return
	}

	if err := h.service.MarkNotificationAsRead(userUUID, notificationUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to mark notification as read")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("notification_id", notificationUUID.String()).Info("Notification marked as read")
	h.SuccessResponse(c, "Notification marked as read", nil)
}

// MarkAllNotificationsAsRead handler untuk mark all notifications as read endpoint
func (h *Handler) MarkAllNotificationsAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.MarkAllNotificationsAsRead(userUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to mark all notifications as read")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("All notifications marked as read")
	h.SuccessResponse(c, "All notifications marked as read", nil)
}

// DeleteNotification handler untuk delete notification endpoint
func (h *Handler) DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	notificationUUID, err := uuid.Parse(c.Param("notificationId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid notification ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", nil)
		return
	}

	if err := h.service.DeleteNotification(userUUID, notificationUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to delete notification")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("notification_id", notificationUUID.String()).Info("Notification deleted successfully")
	h.SuccessResponse(c, "Notification deleted successfully", nil)
}

// GetUserStats handler untuk get user stats endpoint
func (h *Handler) GetUserStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	stats, err := h.service.GetUserStats(userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get user stats")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Info("User stats retrieved successfully")
	h.SuccessResponse(c, "User stats retrieved successfully", stats)
}

// DeactivateAccount handler untuk deactivate account endpoint
func (h *Handler) DeactivateAccount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.DeactivateAccount(userUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to deactivate account")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).Warn("Account deactivated successfully")
	h.SuccessResponse(c, "Account deactivated successfully", nil)
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
