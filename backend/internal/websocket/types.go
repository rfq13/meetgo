package websocket

import (
	"time"

	"github.com/gorilla/websocket"
)

// MessageType mendefinisikan tipe pesan WebSocket
type MessageType string

const (
	// Message types untuk signaling
	MessageTypeOffer        MessageType = "offer"
	MessageTypeAnswer       MessageType = "answer"
	MessageTypeIceCandidate MessageType = "ice-candidate"
	MessageTypeJoinRoom     MessageType = "join-room"
	MessageTypeLeaveRoom    MessageType = "leave-room"
	MessageTypeRoomJoined   MessageType = "room-joined"
	MessageTypeRoomLeft     MessageType = "room-left"
	MessageTypeUserJoined   MessageType = "user-joined"
	MessageTypeUserLeft     MessageType = "user-left"
	MessageTypeError        MessageType = "error"
	MessageTypeSuccess      MessageType = "success"
)

// Message adalah struktur dasar untuk semua pesan WebSocket
type Message struct {
	Type      MessageType `json:"type"`
	RoomID    string      `json:"roomId,omitempty"`
	UserID    string      `json:"userId,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// JoinRoomData adalah data untuk pesan join-room
type JoinRoomData struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

// LeaveRoomData adalah data untuk pesan leave-room
type LeaveRoomData struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

// RoomJoinedData adalah data untuk pesan room-joined
type RoomJoinedData struct {
	RoomID string   `json:"roomId"`
	UserID string   `json:"userId"`
	Users  []string `json:"users"`
}

// RoomLeftData adalah data untuk pesan room-left
type RoomLeftData struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

// UserJoinedData adalah data untuk pesan user-joined
type UserJoinedData struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

// UserLeftData adalah data untuk pesan user-left
type UserLeftData struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
}

// OfferData adalah data untuk pesan offer
type OfferData struct {
	RoomID     string `json:"roomId"`
	FromUserID string `json:"fromUserId"`
	ToUserID   string `json:"toUserId"`
	SDP        string `json:"sdp"`
}

// AnswerData adalah data untuk pesan answer
type AnswerData struct {
	RoomID     string `json:"roomId"`
	FromUserID string `json:"fromUserId"`
	ToUserID   string `json:"toUserId"`
	SDP        string `json:"sdp"`
}

// IceCandidateData adalah data untuk pesan ice-candidate
type IceCandidateData struct {
	RoomID        string `json:"roomId"`
	FromUserID    string `json:"fromUserId"`
	ToUserID      string `json:"toUserId"`
	Candidate     string `json:"candidate"`
	SDPMID        string `json:"sdpMid"`
	SDPMLineIndex int    `json:"sdpMLineIndex"`
}

// ErrorData adalah data untuk pesan error
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessData adalah data untuk pesan success
type SuccessData struct {
	Message string `json:"message"`
}

// Client merepresentasikan sebuah koneksi WebSocket client
type Client struct {
	// Hub adalah pointer ke hub yang mengelola client ini
	Hub *Hub

	// Conn adalah koneksi WebSocket
	Conn *websocket.Conn

	// Send adalah channel untuk mengirim pesan ke client
	Send chan Message

	// UserID adalah ID user yang terhubung
	UserID string

	// RoomIDs adalah daftar room ID yang dijoin oleh client
	RoomIDs map[string]bool
}

// Hub mengelola semua client yang terhubung dan routing pesan
type Hub struct {
	// Clients adalah semua client yang terhubung
	Clients map[*Client]bool

	// Rooms memetakan room ID ke daftar client
	Rooms map[string]map[*Client]bool

	// Register adalah channel untuk registrasi client baru
	Register chan *Client

	// Unregister adalah channel untuk unregister client
	Unregister chan *Client

	// Broadcast adalah channel untuk broadcast pesan ke semua client
	Broadcast chan Message

	// RoomMessage adalah channel untuk mengirim pesan ke room tertentu
	RoomMessage chan RoomMessage

	// DirectMessage adalah channel untuk mengirim pesan langsung ke client tertentu
	DirectMessage chan DirectMessage

	// SignalingHandler untuk WebRTC signaling
	SignalingHandler interface{}
}

// RoomMessage adalah pesan yang akan dikirim ke semua client dalam sebuah room
type RoomMessage struct {
	RoomID  string
	Message Message
}

// DirectMessage adalah pesan yang akan dikirim langsung ke client tertentu
type DirectMessage struct {
	Client  *Client
	Message Message
}

// NewHub membuat instance Hub baru
func NewHub() *Hub {
	return &Hub{
		Clients:       make(map[*Client]bool),
		Rooms:         make(map[string]map[*Client]bool),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan Message),
		RoomMessage:   make(chan RoomMessage),
		DirectMessage: make(chan DirectMessage),
	}
}

// NewClient membuat instance Client baru
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan Message, 256),
		UserID:  userID,
		RoomIDs: make(map[string]bool),
	}
}
