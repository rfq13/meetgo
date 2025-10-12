package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// ReadPump membaca pesan dari koneksi WebSocket
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("error: %v", err)
			}
			break
		}

		// Parse pesan JSON
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			logrus.Errorf("error parsing message: %v", err)
			c.SendErrorMessage("Invalid message format")
			continue
		}

		// Set timestamp
		message.Timestamp = time.Now()

		// Proses pesan berdasarkan tipe
		c.handleMessage(message)
	}
}

// WritePump mengirim pesan ke koneksi WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub menutup channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Serialize pesan ke JSON
			messageBytes, err := json.Marshal(message)
			if err != nil {
				logrus.Errorf("error marshaling message: %v", err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				logrus.Errorf("error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage memproses pesan yang diterima dari client
func (c *Client) handleMessage(message Message) {
	logrus.WithFields(logrus.Fields{
		"type":   message.Type,
		"userId": c.UserID,
		"roomId": message.RoomID,
	}).Info("Handling message")

	switch message.Type {
	case MessageTypeJoinRoom:
		c.handleJoinRoom(message)
	case MessageTypeLeaveRoom:
		c.handleLeaveRoom(message)
	case MessageTypeOffer:
		c.handleOffer(message)
	case MessageTypeAnswer:
		c.handleAnswer(message)
	case MessageTypeIceCandidate:
		c.handleIceCandidate(message)
	default:
		logrus.Warnf("Unknown message type: %s", message.Type)
		c.SendErrorMessage("Unknown message type")
	}
}

// SetSignalingHandler mengatur signaling handler untuk client
func (c *Client) SetSignalingHandler(handler interface{}) {
	// TODO: Implement signaling handler integration
	// Ini akan diimplementasikan setelah signaling handler terintegrasi dengan hub
}

// handleJoinRoom menangani pesan join-room
func (c *Client) handleJoinRoom(message Message) {
	var data JoinRoomData
	if err := mapToStruct(message.Data, &data); err != nil {
		c.SendErrorMessage("Invalid join room data")
		return
	}

	// Validasi data
	if data.RoomID == "" || data.UserID == "" {
		c.SendErrorMessage("Room ID and User ID are required")
		return
	}

	// Pastikan user ID sesuai dengan client
	if data.UserID != c.UserID {
		c.SendErrorMessage("User ID mismatch")
		return
	}

	// Join room melalui hub
	joinMsg := RoomMessage{
		RoomID: data.RoomID,
		Message: Message{
			Type:      MessageTypeJoinRoom,
			RoomID:    data.RoomID,
			UserID:    c.UserID,
			Data:      data,
			Timestamp: time.Now(),
		},
	}

	c.Hub.RoomMessage <- joinMsg
}

// handleLeaveRoom menangani pesan leave-room
func (c *Client) handleLeaveRoom(message Message) {
	var data LeaveRoomData
	if err := mapToStruct(message.Data, &data); err != nil {
		c.SendErrorMessage("Invalid leave room data")
		return
	}

	// Validasi data
	if data.RoomID == "" || data.UserID == "" {
		c.SendErrorMessage("Room ID and User ID are required")
		return
	}

	// Pastikan user ID sesuai dengan client
	if data.UserID != c.UserID {
		c.SendErrorMessage("User ID mismatch")
		return
	}

	// Leave room melalui hub
	leaveMsg := RoomMessage{
		RoomID: data.RoomID,
		Message: Message{
			Type:      MessageTypeLeaveRoom,
			RoomID:    data.RoomID,
			UserID:    c.UserID,
			Data:      data,
			Timestamp: time.Now(),
		},
	}

	c.Hub.RoomMessage <- leaveMsg
}

// handleOffer menangani pesan offer
func (c *Client) handleOffer(message Message) {
	var data OfferData
	if err := mapToStruct(message.Data, &data); err != nil {
		c.SendErrorMessage("Invalid offer data")
		return
	}

	// Validasi data
	if data.RoomID == "" || data.FromUserID == "" || data.ToUserID == "" || data.SDP == "" {
		c.SendErrorMessage("Invalid offer data: missing required fields")
		return
	}

	// Pastikan from user ID sesuai dengan client
	if data.FromUserID != c.UserID {
		c.SendErrorMessage("From user ID mismatch")
		return
	}

	// Kirim offer ke room tertentu
	offerMsg := RoomMessage{
		RoomID: data.RoomID,
		Message: Message{
			Type:      MessageTypeOffer,
			RoomID:    data.RoomID,
			UserID:    c.UserID,
			Data:      data,
			Timestamp: time.Now(),
		},
	}

	c.Hub.RoomMessage <- offerMsg
}

// handleAnswer menangani pesan answer
func (c *Client) handleAnswer(message Message) {
	var data AnswerData
	if err := mapToStruct(message.Data, &data); err != nil {
		c.SendErrorMessage("Invalid answer data")
		return
	}

	// Validasi data
	if data.RoomID == "" || data.FromUserID == "" || data.ToUserID == "" || data.SDP == "" {
		c.SendErrorMessage("Invalid answer data: missing required fields")
		return
	}

	// Pastikan from user ID sesuai dengan client
	if data.FromUserID != c.UserID {
		c.SendErrorMessage("From user ID mismatch")
		return
	}

	// Kirim answer ke room tertentu
	answerMsg := RoomMessage{
		RoomID: data.RoomID,
		Message: Message{
			Type:      MessageTypeAnswer,
			RoomID:    data.RoomID,
			UserID:    c.UserID,
			Data:      data,
			Timestamp: time.Now(),
		},
	}

	c.Hub.RoomMessage <- answerMsg
}

// handleIceCandidate menangani pesan ice-candidate
func (c *Client) handleIceCandidate(message Message) {
	var data IceCandidateData
	if err := mapToStruct(message.Data, &data); err != nil {
		c.SendErrorMessage("Invalid ice candidate data")
		return
	}

	// Validasi data
	if data.RoomID == "" || data.FromUserID == "" || data.ToUserID == "" || data.Candidate == "" {
		c.SendErrorMessage("Invalid ice candidate data: missing required fields")
		return
	}

	// Pastikan from user ID sesuai dengan client
	if data.FromUserID != c.UserID {
		c.SendErrorMessage("From user ID mismatch")
		return
	}

	// Kirim ice candidate ke room tertentu
	iceMsg := RoomMessage{
		RoomID: data.RoomID,
		Message: Message{
			Type:      MessageTypeIceCandidate,
			RoomID:    data.RoomID,
			UserID:    c.UserID,
			Data:      data,
			Timestamp: time.Now(),
		},
	}

	c.Hub.RoomMessage <- iceMsg
}

// SendErrorMessage mengirim pesan error ke client
func (c *Client) SendErrorMessage(message string) {
	errorMsg := Message{
		Type:      MessageTypeError,
		Data:      ErrorData{Code: 400, Message: message},
		Timestamp: time.Now(),
	}

	select {
	case c.Send <- errorMsg:
	default:
		close(c.Send)
	}
}

// SendSuccessMessage mengirim pesan success ke client
func (c *Client) SendSuccessMessage(message string) {
	successMsg := Message{
		Type:      MessageTypeSuccess,
		Data:      SuccessData{Message: message},
		Timestamp: time.Now(),
	}

	select {
	case c.Send <- successMsg:
	default:
		close(c.Send)
	}
}

// mapToStruct mengkonversi interface{} ke struct menggunakan JSON marshaling/unmarshaling
func mapToStruct(data interface{}, result interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, result)
}
