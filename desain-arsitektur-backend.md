# Desain Arsitektur Backend (Golang) untuk WebRTC Meeting

## 1. Arsitektur Umum

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              Load Balancer (Nginx)                            │
└─────────────────────────────────────────────────────────────────────────────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
┌───────────────────▼───────────┐ ┌─────▼────────────────┐ ┌─▼───────────────────┐
│   API Gateway (Gin/Echo)     │ │   WebSocket Server   │ │   Media Server       │
│                              │ │   (Gorilla WS)       │ │   (Pion)             │
│  ┌─────────────────────────┐ │ │                      │ │                      │
│  │ Authentication Service │ │ │  ┌─────────────────┐  │ │  ┌─────────────────┐ │
│  └─────────────────────────┘ │ │  │ Room Manager    │  │ │  │ SFU             │ │
│                             │ │  └─────────────────┘  │ │  └─────────────────┘ │
│  ┌─────────────────────────┐ │ │                      │ │                      │
│  │ User Service            │ │ │  ┌─────────────────┐  │ │  ┌─────────────────┐ │
│  └─────────────────────────┘ │ │  │ Signaling       │  │ │  │ Recording       │ │
│                             │ │  │ Handler         │  │ │  │ Service         │ │
│  ┌─────────────────────────┐ │ │  └─────────────────┘  │ │  └─────────────────┘ │
│  │ Room Service            │ │ │                      │ │                      │
│  └─────────────────────────┘ │ │  ┌─────────────────┐  │ │  ┌─────────────────┐ │
│                             │ │  │ WebRTC Handler  │  │ │  │ Transcoding     │ │
│  ┌─────────────────────────┐ │ │  │   (Pion)        │  │ │  │ Service         │ │
│  │ Notification Service    │ │ │  └─────────────────┘  │ │  └─────────────────┘ │
│  └─────────────────────────┘ │ │                      │ │                      │
└─────────────────────────────┘ └──────────────────────┘ └──────────────────────┘
                    │                    │                    │
                    └────────────────────┼────────────────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    │                    │                    │
        ┌───────────▼──────────┐ ┌──────▼─────────────┐ ┌────▼────────────┐
        │   Database           │ │   Redis Cache      │ │   STUN/TURN     │
        │   (PostgreSQL)       │ │                    │ │   Server        │
        └──────────────────────┘ └───────────────────┘ └─────────────────┘
```

### 1.2 Komponen Utama Backend

#### 1.2.1 API Gateway
- Mengelola HTTP requests untuk REST API
- Mengimplementasikan rate limiting
- Menghandle CORS dan security headers
- Routing ke service yang sesuai

#### 1.2.2 WebSocket Server
- Mengelola koneksi WebSocket real-time
- Menangani signaling untuk WebRTC
- Mengelola room dan participant
- Broadcasting pesan ke client

#### 1.2.3 Media Server
- Mengelola koneksi WebRTC menggunakan Pion
- Implementasi SFU (Selective Forwarding Unit)
- Menangani media processing
- Recording dan transcoding (opsional)

#### 1.2.4 Authentication Service
- Mengelola user registration dan login
- JWT token generation dan validation
- Password hashing dan verification
- Session management

#### 1.2.5 User Service
- Mengelola user profile
- Contact management
- User presence tracking
- User settings management

#### 1.2.6 Room Service
- Room creation dan management
- Room scheduling
- Room history
- Room settings management

#### 1.2.7 Notification Service
- Push notifications
- Email notifications
- In-app notifications
- Webhook notifications

### 1.3 Teknologi Stack

#### 1.3.1 Core Framework
- **Web Framework**: Gin atau Echo untuk HTTP API
- **WebSocket**: Gorilla WebSocket untuk real-time communication
- **WebRTC**: Pion untuk WebRTC implementation
- **Dependency Injection**: Wire atau Dig untuk dependency injection

#### 1.3.2 Database
- **Primary Database**: PostgreSQL dengan GORM
- **Cache**: Redis untuk caching dan session management
- **Search**: Elasticsearch untuk pencarian (opsional)

#### 1.3.3 Message Queue
- **Message Broker**: RabbitMQ atau Kafka untuk asynchronous processing
- **Background Jobs**: Asynq atau Celery untuk background tasks

#### 1.3.4 Monitoring & Logging
- **Logging**: Logrus atau Zap untuk structured logging
- **Metrics**: Prometheus untuk metrics collection
- **Tracing**: Jaeger atau OpenTelemetry untuk distributed tracing
- **Error Tracking**: Sentry untuk error monitoring

## 2. Struktur Folder Backend

```
backend/
├── cmd/                          # Application entry points
│   ├── api/                      # API server
│   │   └── main.go
│   ├── websocket/                # WebSocket server
│   │   └── main.go
│   └── media/                    # Media server
│       └── main.go
├── internal/                     # Private application code
│   ├── api/                      # HTTP handlers
│   │   ├── middleware/           # API middleware
│   │   ├── v1/                   # API version 1
│   │   └── router.go             # API router
│   ├── auth/                     # Authentication service
│   │   ├── handler.go            # Auth handlers
│   │   ├── service.go            # Auth service
│   │   └── repository.go         # Auth repository
│   ├── config/                   # Configuration
│   │   ├── config.go             # Config struct
│   │   └── loader.go             # Config loader
│   ├── database/                 # Database connection
│   │   ├── connection.go         # DB connection
│   │   └── migrations/           # DB migrations
│   ├── models/                   # Database models
│   │   ├── user.go               # User model
│   │   ├── room.go               # Room model
│   │   └── message.go            # Message model
│   ├── notification/             # Notification service
│   │   ├── email.go              # Email notification
│   │   ├── push.go               # Push notification
│   │   └── service.go            # Notification service
│   ├── room/                     # Room service
│   │   ├── handler.go            # Room handlers
│   │   ├── service.go            # Room service
│   │   └── repository.go         # Room repository
│   ├── user/                     # User service
│   │   ├── handler.go            # User handlers
│   │   ├── service.go            # User service
│   │   └── repository.go         # User repository
│   ├── webrtc/                   # WebRTC service
│   │   ├── handler.go            # WebRTC handlers
│   │   ├── peer.go               # WebRTC peer connection
│   │   ├── room.go               # WebRTC room
│   │   └── sfu.go                # SFU implementation
│   ├── websocket/                # WebSocket service
│   │   ├── hub.go                # WebSocket hub
│   │   ├── client.go             # WebSocket client
│   │   └── handler.go            # WebSocket handlers
│   └── utils/                    # Utility functions
│       ├── crypto.go             # Cryptographic utilities
│       ├── validator.go          # Validation utilities
│       └── response.go           # Response utilities
├── pkg/                          # Public library code
│   ├── logger/                   # Logger implementation
│   └── metrics/                  # Metrics implementation
├── scripts/                      # Scripts
│   ├── migration.sh              # Migration script
│   └── deploy.sh                 # Deployment script
├── configs/                      # Configuration files
│   ├── app.yaml                  # Application config
│   ├── database.yaml             # Database config
│   └── redis.yaml                # Redis config
├── deployments/                  # Deployment configurations
│   ├── docker/                   # Docker configs
│   │   ├── Dockerfile            # Dockerfile
│   │   └── docker-compose.yml    # Docker Compose
│   └── kubernetes/               # Kubernetes configs
│       ├── api-deployment.yaml   # API deployment
│       ├── websocket-deployment.yaml # WebSocket deployment
│       └── media-deployment.yaml # Media deployment
├── tests/                        # Test files
│   ├── integration/              # Integration tests
│   ├── unit/                     # Unit tests
│   └── e2e/                      # End-to-end tests
├── go.mod                        # Go modules
├── go.sum                        # Go modules checksum
└── Makefile                      # Build automation
```

## 3. Desain Service Layer

### 3.1 Authentication Service

```go
// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   Repository
	secret string
}

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func NewService(repo Repository, secret string) *Service {
	return &Service{
		repo:   repo,
		secret: secret,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	// Check if user already exists
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *Service) generateJWTToken(user *User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}
```

### 3.2 User Service

```go
// internal/user/service.go
package user

import (
	"context"
	"time"
)

type Service struct {
	repo             Repository
	notificationRepo notification.Repository
}

func NewService(repo Repository, notificationRepo notification.Repository) *Service {
	return &Service{
		repo:             repo,
		notificationRepo: notificationRepo,
	}
}

func (s *Service) GetUserProfile(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Sanitize sensitive data
	user.Password = ""
	return user, nil
}

func (s *Service) UpdateUserProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserContacts(ctx context.Context, userID string) ([]*User, error) {
	return s.repo.FindContacts(ctx, userID)
}

func (s *Service) AddContact(ctx context.Context, userID, contactID string) error {
	return s.repo.AddContact(ctx, userID, contactID)
}

func (s *Service) RemoveContact(ctx context.Context, userID, contactID string) error {
	return s.repo.RemoveContact(ctx, userID, contactID)
}

func (s *Service) UpdateUserStatus(ctx context.Context, userID string, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}

func (s *Service) GetUserPresence(ctx context.Context, userID string) (*UserPresence, error) {
	return s.repo.FindPresence(ctx, userID)
}
```

### 3.3 Room Service

```go
// internal/room/service.go
package room

import (
	"context"
	"errors"
	"time"
)

type Service struct {
	repo             Repository
	userRepo         user.Repository
	notificationRepo notification.Repository
}

func NewService(repo Repository, userRepo user.Repository, notificationRepo notification.Repository) *Service {
	return &Service{
		repo:             repo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *Service) CreateRoom(ctx context.Context, req CreateRoomRequest) (*Room, error) {
	// Validate host exists
	host, err := s.userRepo.FindByID(ctx, req.HostID)
	if err != nil {
		return nil, err
	}
	if host == nil {
		return nil, errors.New("host not found")
	}

	// Create room
	room := &Room{
		Name:        req.Name,
		Description: req.Description,
		HostID:      req.HostID,
		Password:    req.Password,
		MaxUsers:    req.MaxUsers,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *Service) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	return s.repo.FindByID(ctx, roomID)
}

func (s *Service) JoinRoom(ctx context.Context, roomID, userID string, password string) error {
	// Find room
	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return errors.New("room not found")
	}

	// Check room status
	if room.Status != "active" {
		return errors.New("room is not active")
	}

	// Check password if required
	if room.Password != "" && room.Password != password {
		return errors.New("invalid password")
	}

	// Check max users
	participants, err := s.repo.FindParticipants(ctx, roomID)
	if err != nil {
		return err
	}
	if len(participants) >= room.MaxUsers {
		return errors.New("room is full")
	}

	// Join room
	return s.repo.AddParticipant(ctx, roomID, userID)
}

func (s *Service) LeaveRoom(ctx context.Context, roomID, userID string) error {
	return s.repo.RemoveParticipant(ctx, roomID, userID)
}

func (s *Service) GetRoomParticipants(ctx context.Context, roomID string) ([]*User, error) {
	return s.repo.FindParticipants(ctx, roomID)
}

func (s *Service) UpdateRoomSettings(ctx context.Context, roomID string, req UpdateRoomSettingsRequest) (*Room, error) {
	room, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("room not found")
	}

	// Update room settings
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Description != "" {
		room.Description = req.Description
	}
	if req.Password != nil {
		room.Password = *req.Password
	}
	if req.MaxUsers > 0 {
		room.MaxUsers = req.MaxUsers
	}

	room.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *Service) EndRoom(ctx context.Context, roomID string) error {
	return s.repo.UpdateStatus(ctx, roomID, "ended")
}

func (s *Service) GetUserRooms(ctx context.Context, userID string) ([]*Room, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) GetRoomHistory(ctx context.Context, userID string, limit, offset int) ([]*Room, error) {
	return s.repo.FindHistory(ctx, userID, limit, offset)
}
```

### 3.4 WebRTC Service

```go
// internal/webrtc/peer.go
package webrtc

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/pion/webrtc/v3"
)

type Peer struct {
	ID           string
	RoomID       string
	Connection   *webrtc.PeerConnection
	LocalTracks  []*webrtc.TrackLocal
	RemoteTracks map[string]*webrtc.TrackRemote
	DataChannel  *webrtc.DataChannel
	mu           sync.RWMutex
}

type PeerManager struct {
	peers map[string]*Peer
	mu    sync.RWMutex
}

func NewPeerManager() *PeerManager {
	return &PeerManager{
		peers: make(map[string]*Peer),
	}
}

func (pm *PeerManager) CreatePeer(roomID, peerID string) (*Peer, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.peers[peerID]; exists {
		return nil, fmt.Errorf("peer %s already exists", peerID)
	}

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	connection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	peer := &Peer{
		ID:           peerID,
		RoomID:       roomID,
		Connection:   connection,
		RemoteTracks: make(map[string]*webrtc.TrackRemote),
	}

	// Set up event handlers
	connection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		// Send ICE candidate to remote peer via signaling
		candidateJSON := candidate.ToJSON()
		log.Printf("ICE candidate from peer %s: %+v", peerID, candidateJSON)
	})

	connection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("Connection state for peer %s: %s", peerID, state.String())
	})

	connection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Track received from peer %s: %s", peerID, track.Codec().MimeType)
		
		pm.mu.RLock()
		peer.RemoteTracks[track.ID()] = track
		pm.mu.RUnlock()
	})

	connection.OnDataChannel(func(dc *webrtc.DataChannel) {
		log.Printf("Data channel received from peer %s: %s", peerID, dc.Label())
		
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Data channel message from peer %s: %s", peerID, string(msg.Data))
		})
	})

	pm.peers[peerID] = peer
	return peer, nil
}

func (pm *PeerManager) GetPeer(peerID string) (*Peer, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	peer, exists := pm.peers[peerID]
	if !exists {
		return nil, fmt.Errorf("peer %s not found", peerID)
	}

	return peer, nil
}

func (pm *PeerManager) RemovePeer(peerID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	peer, exists := pm.peers[peerID]
	if !exists {
		return fmt.Errorf("peer %s not found", peerID)
	}

	if err := peer.Connection.Close(); err != nil {
		log.Printf("Error closing peer connection for %s: %v", peerID, err)
	}

	delete(pm.peers, peerID)
	return nil
}

func (pm *PeerManager) GetPeersInRoom(roomID string) []*Peer {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var peersInRoom []*Peer
	for _, peer := range pm.peers {
		if peer.RoomID == roomID {
			peersInRoom = append(peersInRoom, peer)
		}
	}

	return peersInRoom
}

func (p *Peer) CreateOffer() (*webrtc.SessionDescription, error) {
	offer, err := p.Connection.CreateOffer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer: %v", err)
	}

	if err := p.Connection.SetLocalDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %v", err)
	}

	return &offer, nil
}

func (p *Peer) CreateAnswer(offer webrtc.SessionDescription) (*webrtc.SessionDescription, error) {
	if err := p.Connection.SetRemoteDescription(offer); err != nil {
		return nil, fmt.Errorf("failed to set remote description: %v", err)
	}

	answer, err := p.Connection.CreateAnswer(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create answer: %v", err)
	}

	if err := p.Connection.SetLocalDescription(answer); err != nil {
		return nil, fmt.Errorf("failed to set local description: %v", err)
	}

	return &answer, nil
}

func (p *Peer) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	return p.Connection.AddICECandidate(candidate)
}

func (p *Peer) AddTrack(track *webrtc.TrackLocal) (*webrtc.RTPSender, error) {
	return p.Connection.AddTrack(track)
}

func (p *Peer) CreateDataChannel(label string) (*webrtc.DataChannel, error) {
	dataChannel, err := p.Connection.CreateDataChannel(label, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create data channel: %v", err)
	}

	p.DataChannel = dataChannel
	return dataChannel, nil
}

func (p *Peer) SendDataMessage(message string) error {
	if p.DataChannel == nil {
		return fmt.Errorf("data channel not available")
	}

	return p.DataChannel.SendText(message)
}

func (p *Peer) Close() error {
	return p.Connection.Close()
}
```

### 3.5 WebSocket Service

```go
// internal/websocket/hub.go
package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	roomID   string
	userID   string
	username string
}

type Room struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	mutex      sync.RWMutex
}

type Hub struct {
	rooms    map[string]*Room
	register chan *Client
	mutex    sync.RWMutex
}

type Message struct {
	Type      string      `json:"type"`
	RoomID    string      `json:"roomId"`
	UserID    string      `json:"userId"`
	Username  string      `json:"username"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		rooms:    make(map[string]*Room),
		register: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleClientRegister(client)

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleClientRegister(client *Client) {
	h.mutex.RLock()
	room, exists := h.rooms[client.roomID]
	h.mutex.RUnlock()

	if !exists {
		room = &Room{
			clients:    make(map[*Client]bool),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan []byte),
		}
		h.mutex.Lock()
		h.rooms[client.roomID] = room
		h.mutex.Unlock()
		go room.Run()
	}

	room.register <- client

	// Send join notification
	joinMessage := Message{
		Type:      "user_joined",
		RoomID:    client.roomID,
		UserID:    client.userID,
		Username:  client.username,
		Payload:   map[string]string{"message": "User joined the room"},
		Timestamp: time.Now(),
	}

	room.broadcast <- h.encodeMessage(joinMessage)
}

func (h *Hub) handleBroadcast(message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	h.mutex.RLock()
	room, exists := h.rooms[msg.RoomID]
	h.mutex.RUnlock()

	if exists {
		room.broadcast <- message
	}
}

func (h *Hub) encodeMessage(msg Message) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return nil
	}
	return data
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.mutex.Lock()
			r.clients[client] = true
			r.mutex.Unlock()

		case client := <-r.unregister:
			r.mutex.Lock()
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			r.mutex.Unlock()

		case message := <-r.broadcast:
			r.mutex.RLock()
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
			r.mutex.RUnlock()
		}
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.mutex.RLock()
		room, exists := hub.rooms[c.roomID]
		hub.mutex.RUnlock()

		if exists {
			room.unregister <- c
		}

		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Set message metadata
		msg.RoomID = c.roomID
		msg.UserID = c.userID
		msg.Username = c.username
		msg.Timestamp = time.Now()

		// Handle message based on type
		switch msg.Type {
		case "chat":
			hub.broadcast <- hub.encodeMessage(msg)

		case "offer", "answer", "ice-candidate":
			// Forward WebRTC signaling messages
			hub.broadcast <- hub.encodeMessage(msg)

		case "user_left":
			// Handle user leaving
			hub.broadcast <- hub.encodeMessage(msg)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
```

## 4. Desain API Endpoints

### 4.1 Authentication Endpoints

```
POST /api/v1/auth/register
Request:
{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe"
}

Response:
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "status": "active",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  }
}

POST /api/v1/auth/login
Request:
{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "status": "active",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    },
    "token": "jwt-token"
  }
}

POST /api/v1/auth/logout
Request:
{
  "token": "jwt-token"
}

Response:
{
  "success": true,
  "message": "Logged out successfully"
}

POST /api/v1/auth/refresh
Request:
{
  "token": "jwt-token"
}

Response:
{
  "success": true,
  "data": {
    "token": "new-jwt-token"
  }
}
```

### 4.2 User Endpoints

```
GET /api/v1/users/profile
Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "avatar": "",
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}

PUT /api/v1/users/profile
Request:
{
  "firstName": "John",
  "lastName": "Doe",
  "avatar": "base64-image",
  "status": "active"
}

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "avatar": "base64-image",
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}

GET /api/v1/users/contacts
Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "email": "contact@example.com",
      "firstName": "Jane",
      "lastName": "Doe",
      "avatar": "",
      "status": "active",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ]
}

POST /api/v1/users/contacts
Request:
{
  "contactId": "uuid"
}

Response:
{
  "success": true,
  "message": "Contact added successfully"
}

DELETE /api/v1/users/contacts/{contactId}
Response:
{
  "success": true,
  "message": "Contact removed successfully"
}
```

### 4.3 Room Endpoints

```
POST /api/v1/rooms
Request:
{
  "name": "Team Meeting",
  "description": "Weekly team sync",
  "password": "password123",
  "maxUsers": 10
}

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "description": "Weekly team sync",
    "hostId": "uuid",
    "password": "password123",
    "maxUsers": 10,
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}

GET /api/v1/rooms/{roomId}
Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "description": "Weekly team sync",
    "hostId": "uuid",
    "password": "password123",
    "maxUsers": 10,
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}

POST /api/v1/rooms/{roomId}/join
Request:
{
  "password": "password123"
}

Response:
{
  "success": true,
  "message": "Joined room successfully"
}

POST /api/v1/rooms/{roomId}/leave
Response:
{
  "success": true,
  "message": "Left room successfully"
}

GET /api/v1/rooms/{roomId}/participants
Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "avatar": "",
      "status": "active"
    }
  ]
}

PUT /api/v1/rooms/{roomId}
Request:
{
  "name": "Team Meeting",
  "description": "Weekly team sync",
  "password": "password123",
  "maxUsers": 15
}

Response:
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Team Meeting",
    "description": "Weekly team sync",
    "hostId": "uuid",
    "password": "password123",
    "maxUsers": 15,
    "status": "active",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}

POST /api/v1/rooms/{roomId}/end
Response:
{
  "success": true,
  "message": "Room ended successfully"
}

GET /api/v1/rooms
Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Team Meeting",
      "description": "Weekly team sync",
      "hostId": "uuid",
      "maxUsers": 10,
      "status": "active",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ]
}

GET /api/v1/rooms/history?page=1&limit=10
Response:
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Team Meeting",
      "description": "Weekly team sync",
      "hostId": "uuid",
      "status": "ended",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1
  }
}
```

### 4.4 WebSocket Endpoints

```
WebSocket: /ws?roomId={roomId}&userId={userId}&token={jwt-token}

Message Types:

1. Join Room
{
  "type": "join_room",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe"
}

2. Chat Message
{
  "type": "chat",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe",
  "payload": {
    "message": "Hello everyone!"
  }
}

3. WebRTC Offer
{
  "type": "offer",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "sdp": "webrtc-sdp-offer"
  }
}

4. WebRTC Answer
{
  "type": "answer",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "sdp": "webrtc-sdp-answer"
  }
}

5. WebRTC ICE Candidate
{
  "type": "ice-candidate",
  "roomId": "uuid",
  "userId": "uuid",
  "targetUserId": "uuid",
  "payload": {
    "candidate": {
      "candidate": "candidate-string",
      "sdpMid": "sdp-mid",
      "sdpMLineIndex": 0
    }
  }
}

6. User Left
{
  "type": "user_left",
  "roomId": "uuid",
  "userId": "uuid",
  "username": "John Doe"
}

7. Mute/Unmute
{
  "type": "mute",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "muted": true
  }
}

8. Video On/Off
{
  "type": "video",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "videoEnabled": false
  }
}

9. Screen Share
{
  "type": "screen_share",
  "roomId": "uuid",
  "userId": "uuid",
  "payload": {
    "enabled": true
  }
}
```

## 5. Kesimpulan

Desain arsitektur backend untuk aplikasi WebRTC meeting ini menggunakan pendekatan microservice dengan pemisahan yang jelas antara komponen-komponen utama:

1. **API Gateway** untuk mengelola HTTP requests dan routing
2. **WebSocket Server** untuk real-time communication dan signaling
3. **Media Server** untuk mengelola koneksi WebRTC menggunakan Pion
4. **Service Layer** untuk business logic (Authentication, User, Room, Notification)
5. **Database Layer** untuk data persistence

Dengan arsitektur ini, aplikasi akan memiliki:
- Scalability yang baik dengan kemampuan untuk menscale komponen secara terpisah
- Maintainability yang tinggi dengan pemisahan tanggung jawab yang jelas
- Performance yang optimal dengan penggunaan teknologi yang tepat untuk setiap komponen
- Reliability yang baik dengan error handling dan monitoring yang memadai

Arsitektur ini juga mendukung fitur-fitur utama yang dibutuhkan seperti video conference real-time, audio conference, room-based meeting system, user management, dan basic UI controls.