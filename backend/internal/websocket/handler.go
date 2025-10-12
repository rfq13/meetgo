package websocket

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Upgrader digunakan untuk mengupgrade HTTP connection ke WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement origin checking untuk production
		// Untuk development, allow all origins
		return true
	},
}

// Handler adalah struct untuk WebSocket HTTP handler
type Handler struct {
	Hub *Hub
}

// NewHandler membuat instance Handler baru
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		Hub: hub,
	}
}

// HandleWebSocket menangani koneksi WebSocket masuk
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// Ambil user ID dari query parameter atau header
	userID := c.Query("userId")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	// Validasi user ID
	if userID == "" {
		logrus.Error("User ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Upgrade HTTP connection ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Log koneksi baru
	logrus.WithFields(logrus.Fields{
		"userId":     userID,
		"remoteAddr": c.Request.RemoteAddr,
	}).Info("New WebSocket connection")

	// Buat client baru
	client := NewClient(h.Hub, conn, userID)

	// Register client ke hub
	h.Hub.Register <- client

	// Start goroutines untuk membaca dan menulis pesan
	go client.WritePump()
	go client.ReadPump()
}

// HandleWebSocketWithAuth menangani koneksi WebSocket dengan autentikasi JWT
func (h *Handler) HandleWebSocketWithAuth(c *gin.Context) {
	// Ambil token dari query parameter atau header
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// Validasi token
	if token == "" {
		logrus.Error("Authorization token is required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// TODO: Validasi JWT token dan dapatkan user ID
	// Untuk saat ini, kita anggap token adalah user ID
	userID := token

	// Upgrade HTTP connection ke WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("Error upgrading connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Log koneksi baru
	logrus.WithFields(logrus.Fields{
		"userId":     userID,
		"remoteAddr": c.Request.RemoteAddr,
	}).Info("New authenticated WebSocket connection")

	// Buat client baru
	client := NewClient(h.Hub, conn, userID)

	// Register client ke hub
	h.Hub.Register <- client

	// Start goroutines untuk membaca dan menulis pesan
	go client.WritePump()
	go client.ReadPump()
}

// GetStats mengembalikan statistik WebSocket server
func (h *Handler) GetStats(c *gin.Context) {
	stats := gin.H{
		"connected_clients": h.Hub.GetClientCount(),
		"active_rooms":      h.Hub.GetRoomCount(),
		"server_status":     "running",
	}

	c.JSON(http.StatusOK, stats)
}

// GetRoomUsers mengembalikan daftar user dalam room tertentu
func (h *Handler) GetRoomUsers(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	users := h.Hub.GetRoomUsers(roomID)

	c.JSON(http.StatusOK, gin.H{
		"roomId": roomID,
		"users":  users,
		"count":  len(users),
	})
}

// GetUserRooms mengembalikan daftar room untuk user tertentu
func (h *Handler) GetUserRooms(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	rooms := h.Hub.GetClientRooms(userID)

	c.JSON(http.StatusOK, gin.H{
		"userId": userID,
		"rooms":  rooms,
		"count":  len(rooms),
	})
}

// BroadcastToRoom mengirim pesan broadcast ke room tertentu (admin endpoint)
func (h *Handler) BroadcastToRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	var request struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat pesan broadcast
	message := Message{
		Type:      MessageType(request.Type),
		RoomID:    roomID,
		Data:      request.Data,
		Timestamp: time.Now(),
	}

	// Kirim ke room
	roomMessage := RoomMessage{
		RoomID:  roomID,
		Message: message,
	}

	h.Hub.RoomMessage <- roomMessage

	c.JSON(http.StatusOK, gin.H{
		"message": "Message broadcasted to room",
		"roomId":  roomID,
		"type":    request.Type,
	})
}

// BroadcastToAll mengirim pesan broadcast ke semua client (admin endpoint)
func (h *Handler) BroadcastToAll(c *gin.Context) {
	var request struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buat pesan broadcast
	message := Message{
		Type:      MessageType(request.Type),
		Data:      request.Data,
		Timestamp: time.Now(),
	}

	// Kirim ke semua client
	h.Hub.Broadcast <- message

	c.JSON(http.StatusOK, gin.H{
		"message": "Message broadcasted to all clients",
		"type":    request.Type,
	})
}

// SendDirectMessage mengirim pesan langsung ke user tertentu (admin endpoint)
func (h *Handler) SendDirectMessage(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var request struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari client target
	var targetClient *Client
	for client := range h.Hub.Clients {
		if client.UserID == userID {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or not connected"})
		return
	}

	// Buat pesan langsung
	message := Message{
		Type:      MessageType(request.Type),
		UserID:    userID,
		Data:      request.Data,
		Timestamp: time.Now(),
	}

	// Kirim ke client target
	directMessage := DirectMessage{
		Client:  targetClient,
		Message: message,
	}

	h.Hub.DirectMessage <- directMessage

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent to user",
		"userId":  userID,
		"type":    request.Type,
	})
}

// DisconnectUser memutuskan koneksi user tertentu (admin endpoint)
func (h *Handler) DisconnectUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Cari client target
	var targetClient *Client
	for client := range h.Hub.Clients {
		if client.UserID == userID {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or not connected"})
		return
	}

	// Tutup koneksi client
	targetClient.Conn.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": "User disconnected",
		"userId":  userID,
	})
}

// SetupRoutes mengatur routes untuk WebSocket handlers
func (h *Handler) SetupRoutes(router *gin.Engine) {
	// WebSocket endpoints
	router.GET("/ws", h.HandleWebSocket)
	router.GET("/ws/auth", h.HandleWebSocketWithAuth)

	// API endpoints untuk monitoring dan admin
	api := router.Group("/api/v1/websocket")
	{
		api.GET("/stats", h.GetStats)
		api.GET("/rooms/:roomId/users", h.GetRoomUsers)
		api.GET("/users/:userId/rooms", h.GetUserRooms)

		// Admin endpoints (TODO: add admin middleware)
		admin := api.Group("/admin")
		{
			admin.POST("/broadcast", h.BroadcastToAll)
			admin.POST("/rooms/:roomId/broadcast", h.BroadcastToRoom)
			admin.POST("/users/:userId/message", h.SendDirectMessage)
			admin.DELETE("/users/:userId/disconnect", h.DisconnectUser)
		}
	}
}
