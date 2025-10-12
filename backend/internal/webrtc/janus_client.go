package webrtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// JanusClient adalah client untuk berkomunikasi dengan Janus WebRTC server
type JanusClient struct {
	// URL Janus server
	BaseURL string

	// Admin URL Janus server (untuk admin API)
	AdminURL string

	// API secret untuk autentikasi
	APISecret string

	// Admin secret untuk admin API
	AdminSecret string

	// HTTP client
	HTTPClient *http.Client

	// Session ID yang aktif
	SessionID uint64

	// Plugin handles yang aktif
	PluginHandles map[uint64]*PluginHandle

	// Mutex untuk thread safety
	mu sync.RWMutex

	// Logger
	logger *logrus.Logger
}

// PluginHandle merepresentasikan handle ke plugin Janus
type PluginHandle struct {
	ID        uint64
	Plugin    string
	SessionID uint64
	Client    *JanusClient
}

// JanusResponse adalah struktur response dari Janus
type JanusResponse struct {
	Janus       string      `json:"janus"`
	Transaction string      `json:"transaction"`
	SessionID   uint64      `json:"session_id,omitempty"`
	HandleID    uint64      `json:"handle_id,omitempty"`
	Plugin      string      `json:"plugin,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Error       *JanusError `json:"error,omitempty"`
}

// JanusError adalah struktur error dari Janus
type JanusError struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

// JanusRequest adalah struktur request ke Janus
type JanusRequest struct {
	Janus       string      `json:"janus"`
	Transaction string      `json:"transaction"`
	SessionID   uint64      `json:"session_id,omitempty"`
	HandleID    uint64      `json:"handle_id,omitempty"`
	Plugin      string      `json:"plugin,omitempty"`
	Body        interface{} `json:"body,omitempty"`
	Jsep        interface{} `json:"jsep,omitempty"`
}

// SessionCreateResponse adalah response untuk session create
type SessionCreateResponse struct {
	Janus       string `json:"janus"`
	Transaction string `json:"transaction"`
	Data        struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

// PluginAttachResponse adalah response untuk plugin attach
type PluginAttachResponse struct {
	Janus       string `json:"janus"`
	Transaction string `json:"transaction"`
	SessionID   uint64 `json:"session_id"`
	Data        struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

// VideoRoomCreateRequest adalah request untuk membuat video room
type VideoRoomCreateRequest struct {
	Request     string   `json:"request"`
	Room        uint64   `json:"room,omitempty"`
	Description string   `json:"description,omitempty"`
	IsPrivate   bool     `json:"is_private,omitempty"`
	Secret      string   `json:"secret,omitempty"`
	Publishers  uint16   `json:"publishers,omitempty"`
	Bitrate     uint64   `json:"bitrate,omitempty"`
	Record      bool     `json:"record,omitempty"`
	RecDir      string   `json:"rec_dir,omitempty"`
	Allowed     []string `json:"allowed,omitempty"`
}

// VideoRoomJoinRequest adalah request untuk join video room
type VideoRoomJoinRequest struct {
	Request string `json:"request"`
	Room    uint64 `json:"room"`
	ID      uint64 `json:"id"`
	Display string `json:"display,omitempty"`
	Token   string `json:"token,omitempty"`
	PTYPE   string `json:"ptype,omitempty"`
}

// VideoRoomPublishRequest adalah request untuk publish ke video room
type VideoRoomPublishRequest struct {
	Request string `json:"request"`
	Audio   bool   `json:"audio,omitempty"`
	Video   bool   `json:"video,omitempty"`
	Data    bool   `json:"data,omitempty"`
}

// VideoRoomSubscribeRequest adalah request untuk subscribe video room
type VideoRoomSubscribeRequest struct {
	Request string                 `json:"request"`
	Room    uint64                 `json:"room"`
	Feed    uint64                 `json:"feed"`
	Jsep    map[string]interface{} `json:"jsep,omitempty"`
}

// JSEP adalah struktur untuk WebRTC session description
type JSEP struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// NewJanusClient membuat instance JanusClient baru
func NewJanusClient(baseURL, adminURL, apiSecret, adminSecret string) *JanusClient {
	return &JanusClient{
		BaseURL:     baseURL,
		AdminURL:    adminURL,
		APISecret:   apiSecret,
		AdminSecret: adminSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		PluginHandles: make(map[uint64]*PluginHandle),
		logger:        logrus.New(),
	}
}

// CreateSession membuat session baru dengan Janus server
func (jc *JanusClient) CreateSession() (uint64, error) {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "create",
		Transaction: transaction,
	}

	if jc.APISecret != "" {
		// TODO: Add API secret to request
	}

	resp, err := jc.makeRequest(request)
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %w", err)
	}

	if resp.Janus != "success" {
		if resp.Error != nil {
			return 0, fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return 0, fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	// Parse session ID dari response
	var sessionResp SessionCreateResponse
	if err := json.Unmarshal(resp.Data.([]byte), &sessionResp); err != nil {
		return 0, fmt.Errorf("failed to parse session response: %w", err)
	}

	jc.mu.Lock()
	jc.SessionID = sessionResp.Data.ID
	jc.mu.Unlock()

	jc.logger.WithField("session_id", jc.SessionID).Info("Created Janus session")

	return jc.SessionID, nil
}

// AttachPlugin men-attach plugin ke session
func (jc *JanusClient) AttachPlugin(pluginName string) (*PluginHandle, error) {
	if jc.SessionID == 0 {
		return nil, fmt.Errorf("no active session")
	}

	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "attach",
		Transaction: transaction,
		SessionID:   jc.SessionID,
		Plugin:      pluginName,
	}

	resp, err := jc.makeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to attach plugin: %w", err)
	}

	if resp.Janus != "success" {
		if resp.Error != nil {
			return nil, fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return nil, fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	// Parse handle ID dari response
	var attachResp PluginAttachResponse
	if err := json.Unmarshal(resp.Data.([]byte), &attachResp); err != nil {
		return nil, fmt.Errorf("failed to parse attach response: %w", err)
	}

	handle := &PluginHandle{
		ID:        attachResp.Data.ID,
		Plugin:    pluginName,
		SessionID: jc.SessionID,
		Client:    jc,
	}

	jc.mu.Lock()
	jc.PluginHandles[handle.ID] = handle
	jc.mu.Unlock()

	jc.logger.WithFields(logrus.Fields{
		"handle_id": handle.ID,
		"plugin":    pluginName,
	}).Info("Attached plugin")

	return handle, nil
}

// makeRequest mengirim request ke Janus server
func (jc *JanusClient) makeRequest(request JanusRequest) (*JanusResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/%s", jc.BaseURL, request.Janus)
	resp, err := jc.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var janusResp JanusResponse
	if err := json.NewDecoder(resp.Body).Decode(&janusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &janusResp, nil
}

// CreateVideoRoom membuat video room baru
func (ph *PluginHandle) CreateVideoRoom(roomID uint64, description string) error {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "message",
		Transaction: transaction,
		SessionID:   ph.SessionID,
		HandleID:    ph.ID,
		Body: VideoRoomCreateRequest{
			Request:     "create",
			Room:        roomID,
			Description: description,
			IsPrivate:   false,
		},
	}

	resp, err := ph.Client.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to create video room: %w", err)
	}

	if resp.Janus != "success" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	ph.Client.logger.WithFields(logrus.Fields{
		"room_id":   roomID,
		"handle_id": ph.ID,
	}).Info("Created video room")

	return nil
}

// JoinVideoRoom bergabung ke video room
func (ph *PluginHandle) JoinVideoRoom(roomID, userID uint64, displayName string) error {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "message",
		Transaction: transaction,
		SessionID:   ph.SessionID,
		HandleID:    ph.ID,
		Body: VideoRoomJoinRequest{
			Request: "join",
			Room:    roomID,
			ID:      userID,
			Display: displayName,
			PTYPE:   "publisher",
		},
	}

	resp, err := ph.Client.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to join video room: %w", err)
	}

	if resp.Janus != "event" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	ph.Client.logger.WithFields(logrus.Fields{
		"room_id":      roomID,
		"user_id":      userID,
		"display_name": displayName,
		"handle_id":    ph.ID,
	}).Info("Joined video room")

	return nil
}

// PublishToVideoRoom mempublish stream ke video room
func (ph *PluginHandle) PublishToVideoRoom(jsep *JSEP) error {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "message",
		Transaction: transaction,
		SessionID:   ph.SessionID,
		HandleID:    ph.ID,
		Body: VideoRoomPublishRequest{
			Request: "publish",
			Audio:   true,
			Video:   true,
		},
	}

	if jsep != nil {
		request.Jsep = jsep
	}

	resp, err := ph.Client.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to publish to video room: %w", err)
	}

	if resp.Janus != "event" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	ph.Client.logger.WithFields(logrus.Fields{
		"handle_id": ph.ID,
	}).Info("Published to video room")

	return nil
}

// SubscribeToVideoRoom subscribe ke publisher di video room
func (ph *PluginHandle) SubscribeToVideoRoom(roomID, feedID uint64, jsep *JSEP) error {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "message",
		Transaction: transaction,
		SessionID:   ph.SessionID,
		HandleID:    ph.ID,
		Body: VideoRoomSubscribeRequest{
			Request: "subscribe",
			Room:    roomID,
			Feed:    feedID,
		},
	}

	if jsep != nil {
		request.Jsep = jsep
	}

	resp, err := ph.Client.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to subscribe to video room: %w", err)
	}

	if resp.Janus != "event" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	ph.Client.logger.WithFields(logrus.Fields{
		"room_id":   roomID,
		"feed_id":   feedID,
		"handle_id": ph.ID,
	}).Info("Subscribed to video room")

	return nil
}

// DestroySession menghapus session Janus
func (jc *JanusClient) DestroySession() error {
	if jc.SessionID == 0 {
		return fmt.Errorf("no active session")
	}

	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "destroy",
		Transaction: transaction,
		SessionID:   jc.SessionID,
	}

	resp, err := jc.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}

	if resp.Janus != "success" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	jc.logger.WithField("session_id", jc.SessionID).Info("Destroyed Janus session")

	jc.mu.Lock()
	jc.SessionID = 0
	jc.PluginHandles = make(map[uint64]*PluginHandle)
	jc.mu.Unlock()

	return nil
}

// DetachPlugin melepas plugin dari session
func (ph *PluginHandle) DetachPlugin() error {
	transaction := uuid.New().String()

	request := JanusRequest{
		Janus:       "detach",
		Transaction: transaction,
		SessionID:   ph.SessionID,
		HandleID:    ph.ID,
	}

	resp, err := ph.Client.makeRequest(request)
	if err != nil {
		return fmt.Errorf("failed to detach plugin: %w", err)
	}

	if resp.Janus != "success" {
		if resp.Error != nil {
			return fmt.Errorf("janus error: %s", resp.Error.Reason)
		}
		return fmt.Errorf("unexpected response: %s", resp.Janus)
	}

	ph.Client.logger.WithFields(logrus.Fields{
		"handle_id": ph.ID,
		"plugin":    ph.Plugin,
	}).Info("Detached plugin")

	ph.Client.mu.Lock()
	delete(ph.Client.PluginHandles, ph.ID)
	ph.Client.mu.Unlock()

	return nil
}
