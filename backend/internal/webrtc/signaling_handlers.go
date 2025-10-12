package webrtc

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/webrtc-meeting/backend/internal/websocket"
)

// SignalingHandler menangani signaling WebRTC dengan integrasi Janus
type SignalingHandler struct {
	// Janus client
	JanusClient *JanusClient

	// Hub WebSocket
	Hub *websocket.Hub

	// Room sessions
	RoomSessions map[string]*RoomSession

	// User sessions
	UserSessions map[string]*UserSession

	// Mutex untuk thread safety
	mu sync.RWMutex

	// Logger
	logger *logrus.Logger
}

// RoomSession merepresentasikan session untuk sebuah room
type RoomSession struct {
	RoomID      string
	JanusRoom   uint64
	Plugin      *PluginHandle
	Publishers  map[string]*PublisherSession
	Subscribers map[string]*SubscriberSession
	CreatedAt   time.Time
}

// PublisherSession merepresentasikan publisher session
type PublisherSession struct {
	UserID       string
	Display      string
	JanusID      uint64
	Plugin       *PluginHandle
	IsPublishing bool
	CreatedAt    time.Time
}

// SubscriberSession merepresentasikan subscriber session
type SubscriberSession struct {
	UserID       string
	FeedID       uint64
	Plugin       *PluginHandle
	IsSubscribed bool
	CreatedAt    time.Time
}

// UserSession merepresentasikan session untuk user
type UserSession struct {
	UserID        string
	RoomIDs       map[string]bool
	PublisherIDs  map[uint64]bool
	SubscriberIDs map[uint64]bool
	CreatedAt     time.Time
}

// NewSignalingHandler membuat instance SignalingHandler baru
func NewSignalingHandler(janusClient *JanusClient, hub *websocket.Hub) *SignalingHandler {
	return &SignalingHandler{
		JanusClient:  janusClient,
		Hub:          hub,
		RoomSessions: make(map[string]*RoomSession),
		UserSessions: make(map[string]*UserSession),
		logger:       logrus.New(),
	}
}

// HandleJoinRoom menangani user yang bergabung ke room
func (sh *SignalingHandler) HandleJoinRoom(roomID, userID, displayName string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.logger.WithFields(logrus.Fields{
		"room_id":      roomID,
		"user_id":      userID,
		"display_name": displayName,
	}).Info("User joining room")

	// Buat atau dapatkan room session
	roomSession, err := sh.getOrCreateRoomSession(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room session: %w", err)
	}

	// Buat atau dapatkan user session
	userSession := sh.getOrCreateUserSession(userID)
	userSession.RoomIDs[roomID] = true

	// Buat publisher session untuk user
	publisherSession := &PublisherSession{
		UserID:    userID,
		Display:   displayName,
		JanusID:   sh.generateJanusUserID(userID),
		CreatedAt: time.Now(),
	}

	// Attach plugin untuk publisher
	publisherPlugin, err := sh.JanusClient.AttachPlugin("janus.plugin.videoroom")
	if err != nil {
		return fmt.Errorf("failed to attach publisher plugin: %w", err)
	}

	publisherSession.Plugin = publisherPlugin
	roomSession.Publishers[userID] = publisherSession
	userSession.PublisherIDs[publisherPlugin.ID] = true

	// Join user ke video room sebagai publisher
	if err := publisherPlugin.JoinVideoRoom(roomSession.JanusRoom, publisherSession.JanusID, displayName); err != nil {
		return fmt.Errorf("failed to join video room: %w", err)
	}

	sh.logger.WithFields(logrus.Fields{
		"room_id":  roomID,
		"user_id":  userID,
		"janus_id": publisherSession.JanusID,
	}).Info("User joined room successfully")

	return nil
}

// HandleLeaveRoom menangani user yang keluar dari room
func (sh *SignalingHandler) HandleLeaveRoom(roomID, userID string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.logger.WithFields(logrus.Fields{
		"room_id": roomID,
		"user_id": userID,
	}).Info("User leaving room")

	// Dapatkan room session
	roomSession, exists := sh.RoomSessions[roomID]
	if !exists {
		return fmt.Errorf("room session not found: %s", roomID)
	}

	// Dapatkan user session
	userSession, exists := sh.UserSessions[userID]
	if !exists {
		return fmt.Errorf("user session not found: %s", userID)
	}

	// Hapus publisher session
	if publisherSession, exists := roomSession.Publishers[userID]; exists {
		if publisherSession.Plugin != nil {
			publisherPlugin := publisherSession.Plugin
			delete(userSession.PublisherIDs, publisherPlugin.ID)

			// Detach plugin
			if err := publisherPlugin.DetachPlugin(); err != nil {
				sh.logger.Errorf("Failed to detach publisher plugin: %v", err)
			}
		}
		delete(roomSession.Publishers, userID)
	}

	// Hapus subscriber sessions untuk user
	for feedID, subscriberSession := range roomSession.Subscribers {
		if subscriberSession.UserID == userID {
			if subscriberSession.Plugin != nil {
				subscriberPlugin := subscriberSession.Plugin
				delete(userSession.SubscriberIDs, subscriberPlugin.ID)

				// Detach plugin
				if err := subscriberPlugin.DetachPlugin(); err != nil {
					sh.logger.Errorf("Failed to detach subscriber plugin: %v", err)
				}
			}
			delete(roomSession.Subscribers, feedID)
		}
	}

	// Update user session
	delete(userSession.RoomIDs, roomID)

	// Jika user tidak ada di room lain, hapus user session
	if len(userSession.RoomIDs) == 0 {
		delete(sh.UserSessions, userID)
	}

	// Jika room tidak ada publisher lagi, hapus room session
	if len(roomSession.Publishers) == 0 && len(roomSession.Subscribers) == 0 {
		delete(sh.RoomSessions, roomID)
	}

	sh.logger.WithFields(logrus.Fields{
		"room_id": roomID,
		"user_id": userID,
	}).Info("User left room successfully")

	return nil
}

// HandleOffer menangani WebRTC offer
func (sh *SignalingHandler) HandleOffer(roomID, fromUserID, toUserID, sdp string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.logger.WithFields(logrus.Fields{
		"room_id":      roomID,
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	}).Info("Handling WebRTC offer")

	// Dapatkan room session
	roomSession, exists := sh.RoomSessions[roomID]
	if !exists {
		return fmt.Errorf("room session not found: %s", roomID)
	}

	// Dapatkan publisher session untuk from user
	publisherSession, exists := roomSession.Publishers[fromUserID]
	if !exists {
		return fmt.Errorf("publisher session not found: %s", fromUserID)
	}

	// Buat JSEP dari offer
	jsep := &JSEP{
		Type: "offer",
		SDP:  sdp,
	}

	// Publish offer ke Janus
	if err := publisherSession.Plugin.PublishToVideoRoom(jsep); err != nil {
		return fmt.Errorf("failed to publish offer: %w", err)
	}

	sh.logger.Info("WebRTC offer handled successfully")

	return nil
}

// HandleAnswer menangani WebRTC answer
func (sh *SignalingHandler) HandleAnswer(roomID, fromUserID, toUserID, sdp string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.logger.WithFields(logrus.Fields{
		"room_id":      roomID,
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
	}).Info("Handling WebRTC answer")

	// Dapatkan room session
	roomSession, exists := sh.RoomSessions[roomID]
	if !exists {
		return fmt.Errorf("room session not found: %s", roomID)
	}

	// Dapatkan publisher session untuk to user (yang akan menerima answer)
	publisherSession, exists := roomSession.Publishers[toUserID]
	if !exists {
		return fmt.Errorf("publisher session not found: %s", toUserID)
	}

	// Buat JSEP dari answer
	jsep := &JSEP{
		Type: "answer",
		SDP:  sdp,
	}

	// Subscribe ke publisher
	if err := publisherSession.Plugin.SubscribeToVideoRoom(roomSession.JanusRoom, publisherSession.JanusID, jsep); err != nil {
		return fmt.Errorf("failed to subscribe with answer: %w", err)
	}

	sh.logger.Info("WebRTC answer handled successfully")

	return nil
}

// HandleIceCandidate menangani WebRTC ICE candidate
func (sh *SignalingHandler) HandleIceCandidate(roomID, fromUserID, toUserID, candidate, sdpMid string, sdpMLineIndex int) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.logger.WithFields(logrus.Fields{
		"room_id":      roomID,
		"from_user_id": fromUserID,
		"to_user_id":   toUserID,
		"candidate":    candidate,
	}).Debug("Handling WebRTC ICE candidate")

	// TODO: Implement ICE candidate handling dengan Janus
	// Untuk saat ini, ICE candidate akan diteruskan melalui WebSocket

	sh.logger.Debug("WebRTC ICE candidate handled successfully")

	return nil
}

// getOrCreateRoomSession membuat atau mendapatkan room session
func (sh *SignalingHandler) getOrCreateRoomSession(roomID string) (*RoomSession, error) {
	if roomSession, exists := sh.RoomSessions[roomID]; exists {
		return roomSession, nil
	}

	// Konversi roomID ke uint64 untuk Janus
	janusRoomID, err := strconv.ParseUint(roomID, 10, 64)
	if err != nil {
		// Jika roomID bukan angka, gunakan hash
		janusRoomID = sh.hashString(roomID)
	}

	// Buat plugin untuk room management
	roomPlugin, err := sh.JanusClient.AttachPlugin("janus.plugin.videoroom")
	if err != nil {
		return nil, fmt.Errorf("failed to attach room plugin: %w", err)
	}

	// Buat room di Janus
	if err := roomPlugin.CreateVideoRoom(janusRoomID, fmt.Sprintf("Room %s", roomID)); err != nil {
		// Room mungkin sudah ada, lanjutkan saja
		sh.logger.Warnf("Failed to create room (might already exist): %v", err)
	}

	// Buat room session
	roomSession := &RoomSession{
		RoomID:      roomID,
		JanusRoom:   janusRoomID,
		Plugin:      roomPlugin,
		Publishers:  make(map[string]*PublisherSession),
		Subscribers: make(map[string]*SubscriberSession),
		CreatedAt:   time.Now(),
	}

	sh.RoomSessions[roomID] = roomSession

	sh.logger.WithFields(logrus.Fields{
		"room_id":       roomID,
		"janus_room_id": janusRoomID,
	}).Info("Created room session")

	return roomSession, nil
}

// getOrCreateUserSession membuat atau mendapatkan user session
func (sh *SignalingHandler) getOrCreateUserSession(userID string) *UserSession {
	if userSession, exists := sh.UserSessions[userID]; exists {
		return userSession
	}

	userSession := &UserSession{
		UserID:        userID,
		RoomIDs:       make(map[string]bool),
		PublisherIDs:  make(map[uint64]bool),
		SubscriberIDs: make(map[uint64]bool),
		CreatedAt:     time.Now(),
	}

	sh.UserSessions[userID] = userSession

	sh.logger.WithField("user_id", userID).Info("Created user session")

	return userSession
}

// generateJanusUserID menghasilkan Janus user ID dari user ID string
func (sh *SignalingHandler) generateJanusUserID(userID string) uint64 {
	// Simple hash function untuk menghasilkan angka dari string
	hash := uint64(0)
	for _, c := range userID {
		hash = hash*31 + uint64(c)
	}

	// Pastikan hash > 0
	if hash == 0 {
		hash = 1
	}

	return hash
}

// hashString fungsi hash sederhana untuk string
func (sh *SignalingHandler) hashString(s string) uint64 {
	hash := uint64(0)
	for _, c := range s {
		hash = hash*31 + uint64(c)
	}

	if hash == 0 {
		hash = 1
	}

	return hash
}

// GetRoomStats mengembalikan statistik room
func (sh *SignalingHandler) GetRoomStats(roomID string) map[string]interface{} {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	stats := map[string]interface{}{
		"room_id": roomID,
		"exists":  false,
	}

	if roomSession, exists := sh.RoomSessions[roomID]; exists {
		stats["exists"] = true
		stats["janus_room_id"] = roomSession.JanusRoom
		stats["publisher_count"] = len(roomSession.Publishers)
		stats["subscriber_count"] = len(roomSession.Subscribers)
		stats["created_at"] = roomSession.CreatedAt

		publishers := make([]string, 0, len(roomSession.Publishers))
		for userID := range roomSession.Publishers {
			publishers = append(publishers, userID)
		}
		stats["publishers"] = publishers

		subscribers := make([]map[string]interface{}, 0, len(roomSession.Subscribers))
		for feedID, subscriber := range roomSession.Subscribers {
			subscribers = append(subscribers, map[string]interface{}{
				"user_id": subscriber.UserID,
				"feed_id": feedID,
			})
		}
		stats["subscribers"] = subscribers
	}

	return stats
}

// GetUserStats mengembalikan statistik user
func (sh *SignalingHandler) GetUserStats(userID string) map[string]interface{} {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	stats := map[string]interface{}{
		"user_id": userID,
		"exists":  false,
	}

	if userSession, exists := sh.UserSessions[userID]; exists {
		stats["exists"] = true
		stats["room_count"] = len(userSession.RoomIDs)
		stats["publisher_count"] = len(userSession.PublisherIDs)
		stats["subscriber_count"] = len(userSession.SubscriberIDs)
		stats["created_at"] = userSession.CreatedAt

		rooms := make([]string, 0, len(userSession.RoomIDs))
		for roomID := range userSession.RoomIDs {
			rooms = append(rooms, roomID)
		}
		stats["rooms"] = rooms
	}

	return stats
}

// GetAllStats mengembalikan statistik semua session
func (sh *SignalingHandler) GetAllStats() map[string]interface{} {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	stats := map[string]interface{}{
		"room_count":   len(sh.RoomSessions),
		"user_count":   len(sh.UserSessions),
		"generated_at": time.Now(),
	}

	roomStats := make(map[string]interface{})
	for roomID := range sh.RoomSessions {
		roomStats[roomID] = sh.GetRoomStats(roomID)
	}
	stats["rooms"] = roomStats

	userStats := make(map[string]interface{})
	for userID := range sh.UserSessions {
		userStats[userID] = sh.GetUserStats(userID)
	}
	stats["users"] = userStats

	return stats
}

// Cleanup melakukan cleanup session yang tidak aktif
func (sh *SignalingHandler) Cleanup(maxAge time.Duration) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	now := time.Now()

	// Cleanup room sessions
	for roomID, roomSession := range sh.RoomSessions {
		if now.Sub(roomSession.CreatedAt) > maxAge {
			sh.logger.WithField("room_id", roomID).Info("Cleaning up expired room session")

			// Detach room plugin
			if roomSession.Plugin != nil {
				if err := roomSession.Plugin.DetachPlugin(); err != nil {
					sh.logger.Errorf("Failed to detach room plugin: %v", err)
				}
			}

			delete(sh.RoomSessions, roomID)
		}
	}

	// Cleanup user sessions
	for userID, userSession := range sh.UserSessions {
		if now.Sub(userSession.CreatedAt) > maxAge && len(userSession.RoomIDs) == 0 {
			sh.logger.WithField("user_id", userID).Info("Cleaning up expired user session")
			delete(sh.UserSessions, userID)
		}
	}
}
