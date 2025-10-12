package websocket

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Run memulai loop utama hub untuk memproses pesan
func (h *Hub) Run() {
	logrus.Info("Starting WebSocket hub")

	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)

		case roomMessage := <-h.RoomMessage:
			h.handleRoomMessage(roomMessage)

		case directMessage := <-h.DirectMessage:
			h.sendDirectMessage(directMessage)
		}
	}
}

// registerClient mendaftarkan client baru ke hub
func (h *Hub) registerClient(client *Client) {
	logrus.WithField("userId", client.UserID).Info("Registering new client")

	h.Clients[client] = true

	// Kirim pesan success ke client
	client.SendSuccessMessage("Connected to WebSocket server")
}

// unregisterClient menghapus client dari hub dan semua room
func (h *Hub) unregisterClient(client *Client) {
	logrus.WithField("userId", client.UserID).Info("Unregistering client")

	if _, ok := h.Clients[client]; ok {
		delete(h.Clients, client)

		// Hapus client dari semua room
		for roomID := range client.RoomIDs {
			h.leaveRoom(client, roomID)
		}

		close(client.Send)
	}
}

// broadcastMessage mengirim pesan ke semua client yang terhubung
func (h *Hub) broadcastMessage(message Message) {
	logrus.WithFields(logrus.Fields{
		"type": message.Type,
	}).Info("Broadcasting message")

	for client := range h.Clients {
		select {
		case client.Send <- message:
		default:
			// Client tidak bisa menerima pesan, unregister
			close(client.Send)
			delete(h.Clients, client)
		}
	}
}

// handleRoomMessage memproses pesan yang ditujukan untuk room tertentu
func (h *Hub) handleRoomMessage(roomMessage RoomMessage) {
	message := roomMessage.Message
	roomID := roomMessage.RoomID

	logrus.WithFields(logrus.Fields{
		"type":   message.Type,
		"roomId": roomID,
		"userId": message.UserID,
	}).Info("Handling room message")

	switch message.Type {
	case MessageTypeJoinRoom:
		h.handleJoinRoom(message, roomID)
	case MessageTypeLeaveRoom:
		h.handleLeaveRoom(message, roomID)
	case MessageTypeOffer:
		h.handleOfferMessage(message, roomID)
	case MessageTypeAnswer:
		h.handleAnswerMessage(message, roomID)
	case MessageTypeIceCandidate:
		h.handleIceCandidateMessage(message, roomID)
	default:
		logrus.Warnf("Unknown room message type: %s", message.Type)
	}
}

// handleJoinRoom menangani client yang bergabung ke room
func (h *Hub) handleJoinRoom(message Message, roomID string) {
	var data JoinRoomData
	if err := mapToStruct(message.Data, &data); err != nil {
		logrus.Errorf("Error parsing join room data: %v", err)
		return
	}

	// Gunakan signaling handler jika tersedia
	if h.SignalingHandler != nil {
		if signalingHandler, ok := h.SignalingHandler.(interface {
			HandleJoinRoom(roomID, userID, displayName string) error
		}); ok {
			if err := signalingHandler.HandleJoinRoom(roomID, data.UserID, data.UserID); err != nil {
				logrus.Errorf("Error handling join room with signaling handler: %v", err)
			}
		}
	}

	// Cari client yang mengirim pesan
	var sender *Client
	for client := range h.Clients {
		if client.UserID == data.UserID {
			sender = client
			break
		}
	}

	if sender == nil {
		logrus.Error("Sender client not found")
		return
	}

	// Tambahkan client ke room
	if _, exists := h.Rooms[roomID]; !exists {
		h.Rooms[roomID] = make(map[*Client]bool)
	}

	// Jika client sudah ada di room, tidak perlu ditambahkan lagi
	if !sender.RoomIDs[roomID] {
		h.Rooms[roomID][sender] = true
		sender.RoomIDs[roomID] = true

		logrus.WithFields(logrus.Fields{
			"roomId": roomID,
			"userId": sender.UserID,
		}).Info("User joined room")

		// Kumpulkan daftar user di room
		var users []string
		for client := range h.Rooms[roomID] {
			users = append(users, client.UserID)
		}

		// Kirim konfirmasi ke sender
		roomJoinedMsg := Message{
			Type:   MessageTypeRoomJoined,
			RoomID: roomID,
			UserID: sender.UserID,
			Data: RoomJoinedData{
				RoomID: roomID,
				UserID: sender.UserID,
				Users:  users,
			},
			Timestamp: time.Now(),
		}

		select {
		case sender.Send <- roomJoinedMsg:
		default:
			close(sender.Send)
		}

		// Beritahu user lain di room bahwa ada user baru bergabung
		userJoinedMsg := Message{
			Type:   MessageTypeUserJoined,
			RoomID: roomID,
			UserID: sender.UserID,
			Data: UserJoinedData{
				RoomID: roomID,
				UserID: sender.UserID,
			},
			Timestamp: time.Now(),
		}

		h.broadcastToRoom(roomID, userJoinedMsg, sender)
	} else {
		// Client sudah ada di room, kirim daftar user saat ini
		var users []string
		for client := range h.Rooms[roomID] {
			users = append(users, client.UserID)
		}

		roomJoinedMsg := Message{
			Type:   MessageTypeRoomJoined,
			RoomID: roomID,
			UserID: sender.UserID,
			Data: RoomJoinedData{
				RoomID: roomID,
				UserID: sender.UserID,
				Users:  users,
			},
			Timestamp: time.Now(),
		}

		select {
		case sender.Send <- roomJoinedMsg:
		default:
			close(sender.Send)
		}
	}
}

// handleLeaveRoom menangani client yang keluar dari room
func (h *Hub) handleLeaveRoom(message Message, roomID string) {
	var data LeaveRoomData
	if err := mapToStruct(message.Data, &data); err != nil {
		logrus.Errorf("Error parsing leave room data: %v", err)
		return
	}

	// Gunakan signaling handler jika tersedia
	if h.SignalingHandler != nil {
		if signalingHandler, ok := h.SignalingHandler.(interface {
			HandleLeaveRoom(roomID, userID string) error
		}); ok {
			if err := signalingHandler.HandleLeaveRoom(roomID, data.UserID); err != nil {
				logrus.Errorf("Error handling leave room with signaling handler: %v", err)
			}
		}
	}

	// Cari client yang mengirim pesan
	var sender *Client
	for client := range h.Clients {
		if client.UserID == data.UserID {
			sender = client
			break
		}
	}

	if sender == nil {
		logrus.Error("Sender client not found")
		return
	}

	// Hapus client dari room
	h.leaveRoom(sender, roomID)

	// Kirim konfirmasi ke sender
	roomLeftMsg := Message{
		Type:   MessageTypeRoomLeft,
		RoomID: roomID,
		UserID: sender.UserID,
		Data: RoomLeftData{
			RoomID: roomID,
			UserID: sender.UserID,
		},
		Timestamp: time.Now(),
	}

	select {
	case sender.Send <- roomLeftMsg:
	default:
		close(sender.Send)
	}

	// Beritahu user lain di room bahwa user telah keluar
	userLeftMsg := Message{
		Type:   MessageTypeUserLeft,
		RoomID: roomID,
		UserID: sender.UserID,
		Data: UserLeftData{
			RoomID: roomID,
			UserID: sender.UserID,
		},
		Timestamp: time.Now(),
	}

	h.broadcastToRoom(roomID, userLeftMsg, sender)
}

// leaveRoom menghapus client dari room
func (h *Hub) leaveRoom(client *Client, roomID string) {
	if roomClients, exists := h.Rooms[roomID]; exists {
		if roomClients[client] {
			delete(roomClients, client)
			delete(client.RoomIDs, roomID)

			// Jika room kosong, hapus room
			if len(roomClients) == 0 {
				delete(h.Rooms, roomID)
			}

			logrus.WithFields(logrus.Fields{
				"roomId": roomID,
				"userId": client.UserID,
			}).Info("User left room")
		}
	}
}

// handleOfferMessage menangani pesan offer
func (h *Hub) handleOfferMessage(message Message, roomID string) {
	var data OfferData
	if err := mapToStruct(message.Data, &data); err != nil {
		logrus.Errorf("Error parsing offer data: %v", err)
		return
	}

	// Gunakan signaling handler jika tersedia
	if h.SignalingHandler != nil {
		if signalingHandler, ok := h.SignalingHandler.(interface {
			HandleOffer(roomID, fromUserID, toUserID, sdp string) error
		}); ok {
			if err := signalingHandler.HandleOffer(roomID, data.FromUserID, data.ToUserID, data.SDP); err != nil {
				logrus.Errorf("Error handling offer with signaling handler: %v", err)
			}
			return
		}
	}

	// Fallback ke routing langsung jika tidak ada signaling handler
	// Cari client target
	var targetClient *Client
	for client := range h.Clients {
		if client.UserID == data.ToUserID {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		logrus.Warnf("Target client not found: %s", data.ToUserID)
		return
	}

	// Pastikan target client ada di room yang sama
	if !targetClient.RoomIDs[roomID] {
		logrus.Warnf("Target client not in room: %s", roomID)
		return
	}

	// Kirim offer ke target client
	select {
	case targetClient.Send <- message:
	default:
		close(targetClient.Send)
		delete(h.Clients, targetClient)
	}
}

// handleAnswerMessage menangani pesan answer
func (h *Hub) handleAnswerMessage(message Message, roomID string) {
	var data AnswerData
	if err := mapToStruct(message.Data, &data); err != nil {
		logrus.Errorf("Error parsing answer data: %v", err)
		return
	}

	// Gunakan signaling handler jika tersedia
	if h.SignalingHandler != nil {
		if signalingHandler, ok := h.SignalingHandler.(interface {
			HandleAnswer(roomID, fromUserID, toUserID, sdp string) error
		}); ok {
			if err := signalingHandler.HandleAnswer(roomID, data.FromUserID, data.ToUserID, data.SDP); err != nil {
				logrus.Errorf("Error handling answer with signaling handler: %v", err)
			}
			return
		}
	}

	// Fallback ke routing langsung jika tidak ada signaling handler
	// Cari client target
	var targetClient *Client
	for client := range h.Clients {
		if client.UserID == data.ToUserID {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		logrus.Warnf("Target client not found: %s", data.ToUserID)
		return
	}

	// Pastikan target client ada di room yang sama
	if !targetClient.RoomIDs[roomID] {
		logrus.Warnf("Target client not in room: %s", roomID)
		return
	}

	// Kirim answer ke target client
	select {
	case targetClient.Send <- message:
	default:
		close(targetClient.Send)
		delete(h.Clients, targetClient)
	}
}

// handleIceCandidateMessage menangani pesan ice-candidate
func (h *Hub) handleIceCandidateMessage(message Message, roomID string) {
	var data IceCandidateData
	if err := mapToStruct(message.Data, &data); err != nil {
		logrus.Errorf("Error parsing ice candidate data: %v", err)
		return
	}

	// Gunakan signaling handler jika tersedia
	if h.SignalingHandler != nil {
		if signalingHandler, ok := h.SignalingHandler.(interface {
			HandleIceCandidate(roomID, fromUserID, toUserID, candidate, sdpMid string, sdpMLineIndex int) error
		}); ok {
			if err := signalingHandler.HandleIceCandidate(roomID, data.FromUserID, data.ToUserID, data.Candidate, data.SDPMID, data.SDPMLineIndex); err != nil {
				logrus.Errorf("Error handling ice candidate with signaling handler: %v", err)
			}
			return
		}
	}

	// Fallback ke routing langsung jika tidak ada signaling handler
	// Cari client target
	var targetClient *Client
	for client := range h.Clients {
		if client.UserID == data.ToUserID {
			targetClient = client
			break
		}
	}

	if targetClient == nil {
		logrus.Warnf("Target client not found: %s", data.ToUserID)
		return
	}

	// Pastikan target client ada di room yang sama
	if !targetClient.RoomIDs[roomID] {
		logrus.Warnf("Target client not in room: %s", roomID)
		return
	}

	// Kirim ice candidate ke target client
	select {
	case targetClient.Send <- message:
	default:
		close(targetClient.Send)
		delete(h.Clients, targetClient)
	}
}

// broadcastToRoom mengirim pesan ke semua client dalam room kecuali sender
func (h *Hub) broadcastToRoom(roomID string, message Message, sender *Client) {
	if roomClients, exists := h.Rooms[roomID]; exists {
		for client := range roomClients {
			if client != sender {
				select {
				case client.Send <- message:
				default:
					// Client tidak bisa menerima pesan, unregister
					close(client.Send)
					delete(h.Clients, client)
					delete(roomClients, client)
				}
			}
		}

		// Jika room kosong, hapus room
		if len(roomClients) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}

// sendDirectMessage mengirim pesan langsung ke client tertentu
func (h *Hub) sendDirectMessage(directMessage DirectMessage) {
	client := directMessage.Client
	message := directMessage.Message

	logrus.WithFields(logrus.Fields{
		"type":   message.Type,
		"userId": client.UserID,
	}).Info("Sending direct message")

	select {
	case client.Send <- message:
	default:
		// Client tidak bisa menerima pesan, unregister
		close(client.Send)
		delete(h.Clients, client)
	}
}

// GetRoomUsers mengembalikan daftar user ID dalam room tertentu
func (h *Hub) GetRoomUsers(roomID string) []string {
	var users []string

	if roomClients, exists := h.Rooms[roomID]; exists {
		for client := range roomClients {
			users = append(users, client.UserID)
		}
	}

	return users
}

// GetClientRooms mengembalikan daftar room ID untuk client tertentu
func (h *Hub) GetClientRooms(userID string) []string {
	var rooms []string

	for client := range h.Clients {
		if client.UserID == userID {
			for roomID := range client.RoomIDs {
				rooms = append(rooms, roomID)
			}
			break
		}
	}

	return rooms
}

// IsClientInRoom memeriksa apakah client ada di room tertentu
func (h *Hub) IsClientInRoom(userID, roomID string) bool {
	for client := range h.Clients {
		if client.UserID == userID {
			return client.RoomIDs[roomID]
		}
	}
	return false
}

// GetClientCount mengembalikan jumlah client yang terhubung
func (h *Hub) GetClientCount() int {
	return len(h.Clients)
}

// GetRoomCount mengembalikan jumlah room yang aktif
func (h *Hub) GetRoomCount() int {
	return len(h.Rooms)
}
