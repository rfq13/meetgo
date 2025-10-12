package room

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Handler struct untuk room handler
type Handler struct {
	service                *Service
	logger                 *logger.Logger
	authMiddleware         gin.HandlerFunc
	optionalAuthMiddleware gin.HandlerFunc
}

// NewHandler membuat room handler baru
func NewHandler(service *Service, log *logger.Logger) *Handler {
	return &Handler{
		service:                service,
		logger:                 log,
		authMiddleware:         func(c *gin.Context) { c.Next() },
		optionalAuthMiddleware: func(c *gin.Context) { c.Next() },
	}
}

// RegisterRoutes registrasi routes untuk room management
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	rooms := router.Group("/rooms")
	{
		// Public routes (with optional auth)
		rooms.GET("", h.OptionalAuthMiddleware(), h.GetRooms)
		rooms.GET("/:roomId", h.OptionalAuthMiddleware(), h.GetRoom)

		// Protected routes
		rooms.POST("", h.AuthMiddleware(), h.CreateRoom)
		rooms.PUT("/:roomId", h.AuthMiddleware(), h.UpdateRoom)
		rooms.DELETE("/:roomId", h.AuthMiddleware(), h.DeleteRoom)

		// Room participation
		rooms.POST("/:roomId/join", h.AuthMiddleware(), h.JoinRoom)
		rooms.POST("/:roomId/leave", h.AuthMiddleware(), h.LeaveRoom)
		rooms.POST("/:roomId/end", h.AuthMiddleware(), h.EndRoom)

		// Room participants
		rooms.GET("/:roomId/participants", h.AuthMiddleware(), h.GetRoomParticipants)
		rooms.POST("/:roomId/participants/:participantId/kick", h.AuthMiddleware(), h.KickParticipant)

		// Room messages
		rooms.GET("/:roomId/messages", h.AuthMiddleware(), h.GetRoomMessages)

		// Room settings
		rooms.GET("/:roomId/settings", h.AuthMiddleware(), h.GetRoomSettings)
		rooms.PUT("/:roomId/settings", h.AuthMiddleware(), h.UpdateRoomSettings)

		// Room stats
		rooms.GET("/:roomId/stats", h.AuthMiddleware(), h.GetRoomStats)
	}
}

// CreateRoom handler untuk create room endpoint
func (h *Handler) CreateRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid create room request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	room, err := h.service.CreateRoom(userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to create room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", room.ID.String()).Info("Room created successfully")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Room created successfully",
		"data":    room,
	})
}

// GetRooms handler untuk get rooms endpoint
func (h *Handler) GetRooms(c *gin.Context) {
	params := h.GetPaginationParams(c)
	status := c.Query("status")
	hostID := c.Query("host_id")

	// Get user ID if authenticated
	var userID uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			if parsedUID, err := uuid.Parse(uidStr); err == nil {
				userID = parsedUID
			}
		}
	}

	rooms, total, err := h.service.GetRooms(params.Page, params.PerPage, status, hostID, userID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get rooms")
		h.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	h.logger.WithField("total", total).Info("Rooms retrieved successfully")
	h.PaginatedResponse(c, "Rooms retrieved successfully", rooms, total, params)
}

// GetRoom handler untuk get room endpoint
func (h *Handler) GetRoom(c *gin.Context) {
	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	// Get user ID if authenticated
	var userID uuid.UUID
	if uid, exists := c.Get("user_id"); exists {
		if uidStr, ok := uid.(string); ok {
			if parsedUID, err := uuid.Parse(uidStr); err == nil {
				userID = parsedUID
			}
		}
	}

	room, err := h.service.GetRoom(roomUUID, userID)
	if err != nil {
		h.logger.WithError(err).WithField("room_id", roomUUID).Error("Failed to get room")
		h.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	h.logger.WithField("room_id", roomUUID.String()).Info("Room retrieved successfully")
	h.SuccessResponse(c, "Room retrieved successfully", room)
}

// UpdateRoom handler untuk update room endpoint
func (h *Handler) UpdateRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update room request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	room, err := h.service.UpdateRoom(roomUUID, userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to update room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room updated successfully")
	h.SuccessResponse(c, "Room updated successfully", room)
}

// DeleteRoom handler untuk delete room endpoint
func (h *Handler) DeleteRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.DeleteRoom(roomUUID, userUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to delete room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room deleted successfully")
	h.SuccessResponse(c, "Room deleted successfully", nil)
}

// JoinRoom handler untuk join room endpoint
func (h *Handler) JoinRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	var req JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid join room request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	participant, err := h.service.JoinRoom(roomUUID, userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to join room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("User joined room successfully")
	h.SuccessResponse(c, "Joined room successfully", participant)
}

// LeaveRoom handler untuk leave room endpoint
func (h *Handler) LeaveRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.LeaveRoom(roomUUID, userUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to leave room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("User left room successfully")
	h.SuccessResponse(c, "Left room successfully", nil)
}

// EndRoom handler untuk end room endpoint
func (h *Handler) EndRoom(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.EndRoom(roomUUID, userUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to end room")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room ended successfully")
	h.SuccessResponse(c, "Room ended successfully", nil)
}

// GetRoomParticipants handler untuk get room participants endpoint
func (h *Handler) GetRoomParticipants(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	participants, err := h.service.GetRoomParticipants(roomUUID, userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get room participants")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room participants retrieved successfully")
	h.SuccessResponse(c, "Room participants retrieved successfully", participants)
}

// KickParticipant handler untuk kick participant endpoint
func (h *Handler) KickParticipant(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	participantUUID, err := uuid.Parse(c.Param("participantId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid participant ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid participant ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	if err := h.service.KickParticipant(roomUUID, userUUID, participantUUID); err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to kick participant")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).WithField("participant_id", participantUUID.String()).Info("Participant kicked successfully")
	h.SuccessResponse(c, "Participant kicked successfully", nil)
}

// GetRoomMessages handler untuk get room messages endpoint
func (h *Handler) GetRoomMessages(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	params := h.GetPaginationParams(c)

	messages, total, err := h.service.GetRoomMessages(roomUUID, userUUID, params.Page, params.PerPage)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get room messages")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room messages retrieved successfully")
	h.PaginatedResponse(c, "Room messages retrieved successfully", messages, total, params)
}

// GetRoomSettings handler untuk get room settings endpoint
func (h *Handler) GetRoomSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	settings, err := h.service.GetRoomSettings(roomUUID, userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get room settings")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room settings retrieved successfully")
	h.SuccessResponse(c, "Room settings retrieved successfully", settings)
}

// UpdateRoomSettings handler untuk update room settings endpoint
func (h *Handler) UpdateRoomSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	var req UpdateRoomSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update room settings request")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	settings, err := h.service.UpdateRoomSettings(roomUUID, userUUID, &req)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to update room settings")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room settings updated successfully")
	h.SuccessResponse(c, "Room settings updated successfully", settings)
}

// GetRoomStats handler untuk get room stats endpoint
func (h *Handler) GetRoomStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomUUID, err := uuid.Parse(c.Param("roomId"))
	if err != nil {
		h.logger.WithError(err).Error("Invalid room ID in parameter")
		h.ErrorResponse(c, http.StatusBadRequest, "Invalid room ID", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID in context")
		h.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	stats, err := h.service.GetRoomStats(roomUUID, userUUID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userUUID).Error("Failed to get room stats")
		h.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	h.logger.WithUserID(userUUID.String()).WithField("room_id", roomUUID.String()).Info("Room stats retrieved successfully")
	h.SuccessResponse(c, "Room stats retrieved successfully", stats)
}

// AuthMiddleware middleware untuk authentication (akan menggunakan middleware yang diinject)
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	if h.authMiddleware != nil {
		return h.authMiddleware
	}
	return func(c *gin.Context) { c.Next() }
}

// OptionalAuthMiddleware middleware untuk optional authentication (akan menggunakan middleware yang diinject)
func (h *Handler) OptionalAuthMiddleware() gin.HandlerFunc {
	if h.optionalAuthMiddleware != nil {
		return h.optionalAuthMiddleware
	}
	return func(c *gin.Context) { c.Next() }
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

// SetAuthMiddleware sets the actual auth middleware
func (h *Handler) SetAuthMiddleware(middleware gin.HandlerFunc) {
	if middleware != nil {
		h.authMiddleware = middleware
	} else {
		h.authMiddleware = func(c *gin.Context) { c.Next() }
	}
}

// SetOptionalAuthMiddleware sets the actual optional auth middleware
func (h *Handler) SetOptionalAuthMiddleware(middleware gin.HandlerFunc) {
	if middleware != nil {
		h.optionalAuthMiddleware = middleware
	} else {
		h.optionalAuthMiddleware = func(c *gin.Context) { c.Next() }
	}
}
