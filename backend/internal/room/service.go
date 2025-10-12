package room

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/webrtc-meeting/backend/models"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Service struct untuk room service
type Service struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewService membuat room service baru
func NewService(db *gorm.DB, log *logger.Logger) *Service {
	return &Service{
		db:     db,
		logger: log,
	}
}

// CreateRoomRequest struct untuk request create room
type CreateRoomRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=100"`
	Description string     `json:"description" binding:"max=500"`
	Password    string     `json:"password" binding:"min=6"`
	MaxUsers    int        `json:"max_users" binding:"min=2,max=100"`
	Type        string     `json:"type" binding:"required,oneof=meeting webinar conference classroom"`
	IsPublic    bool       `json:"is_public"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
}

// UpdateRoomRequest struct untuk request update room
type UpdateRoomRequest struct {
	Name        string     `json:"name" binding:"required,min=1,max=100"`
	Description string     `json:"description" binding:"max=500"`
	Password    string     `json:"password" binding:"omitempty,min=6"`
	MaxUsers    int        `json:"max_users" binding:"min=2,max=100"`
	IsPublic    bool       `json:"is_public"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
}

// JoinRoomRequest struct untuk request join room
type JoinRoomRequest struct {
	Password string `json:"password"`
}

// UpdateRoomSettingsRequest struct untuk request update room settings
type UpdateRoomSettingsRequest struct {
	AllowScreenShare    bool   `json:"allow_screen_share"`
	AllowChat           bool   `json:"allow_chat"`
	AllowFileShare      bool   `json:"allow_file_share"`
	RequirePassword     bool   `json:"require_password"`
	WaitingRoom         bool   `json:"waiting_room"`
	AutoRecord          bool   `json:"auto_record"`
	MaxParticipants     int    `json:"max_participants" binding:"min=2,max=100"`
	VideoQuality        string `json:"video_quality"`
	AudioQuality        string `json:"audio_quality"`
	EnableBreakoutRooms bool   `json:"enable_breakout_rooms"`
	EnablePolling       bool   `json:"enable_polling"`
	EnableWhiteboard    bool   `json:"enable_whiteboard"`
	EnableRecording     bool   `json:"enable_recording"`
}

// CreateRoom membuat room baru
func (s *Service) CreateRoom(userID uuid.UUID, req *CreateRoomRequest) (*models.Room, error) {
	// Validate room type
	roomType := models.RoomType(req.Type)
	if !isValidRoomType(roomType) {
		return nil, fmt.Errorf("invalid room type")
	}

	// Validate time range
	if req.StartTime != nil && req.EndTime != nil {
		if req.EndTime.Before(*req.StartTime) {
			return nil, fmt.Errorf("end time must be after start time")
		}
	}

	// Generate unique room code
	roomCode, err := s.generateRoomCode()
	if err != nil {
		s.logger.LogError(err, "Failed to generate room code")
		return nil, fmt.Errorf("failed to generate room code")
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.LogError(err, "Failed to hash room password")
			return nil, fmt.Errorf("failed to process password")
		}
		hashedPassword = string(hashedBytes)
	}

	// Create room
	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
		HostID:      userID,
		RoomCode:    roomCode,
		Password:    hashedPassword,
		MaxUsers:    req.MaxUsers,
		Status:      models.RoomStatusActive,
		Type:        roomType,
		IsPublic:    req.IsPublic,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	if err := s.db.Create(room).Error; err != nil {
		s.logger.LogError(err, "Failed to create room")
		return nil, fmt.Errorf("failed to create room")
	}

	// Create default room settings
	roomSettings := &models.RoomSetting{
		RoomID:              room.ID,
		AllowScreenShare:    true,
		AllowChat:           true,
		AllowFileShare:      true,
		RequirePassword:     req.Password != "",
		WaitingRoom:         false,
		AutoRecord:          false,
		MaxParticipants:     req.MaxUsers,
		VideoQuality:        "hd",
		AudioQuality:        "high",
		EnableBreakoutRooms: false,
		EnablePolling:       false,
		EnableWhiteboard:    false,
		EnableRecording:     true,
	}

	if err := s.db.Create(roomSettings).Error; err != nil {
		s.logger.LogError(err, "Failed to create room settings")
	}

	// Load room with relations
	if err := s.db.Preload("Host").Preload("Settings").First(room, room.ID).Error; err != nil {
		s.logger.LogError(err, "Failed to load room with relations")
		return nil, fmt.Errorf("failed to load room")
	}

	// Clear password before returning
	room.Password = ""

	s.logger.WithUserID(userID.String()).WithField("room_id", room.ID.String()).Info("Room created successfully")
	return room, nil
}

// GetRooms mengambil daftar rooms
func (s *Service) GetRooms(page, perPage int, status, hostID string, userID uuid.UUID) ([]*models.Room, int64, error) {
	var rooms []*models.Room
	var total int64

	query := s.db.Model(&models.Room{})

	// Add filters
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if hostID != "" {
		if hostUUID, err := uuid.Parse(hostID); err == nil {
			query = query.Where("host_id = ?", hostUUID)
		}
	}

	// For non-public rooms, only show rooms where user is host or participant
	query = query.Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
		true, userID, userID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		s.logger.LogError(err, "Failed to count rooms")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Get rooms with pagination and relations
	offset := (page - 1) * perPage
	if err := query.Preload("Host").
		Preload("Participants").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&rooms).Error; err != nil {
		s.logger.LogError(err, "Failed to get rooms")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Clear passwords
	for _, room := range rooms {
		room.Password = ""
	}

	return rooms, total, nil
}

// GetRoom mengambil detail room
func (s *Service) GetRoom(roomID uuid.UUID, userID uuid.UUID) (*models.Room, error) {
	var room models.Room

	// Check if user has access to room (host, participant, or public room)
	query := s.db.Where("id = ?", roomID).
		Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
			true, userID, userID)

	if err := query.Preload("Host").
		Preload("Participants").
		Preload("Participants.User").
		Preload("Settings").
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to get room")
		return nil, fmt.Errorf("internal server error")
	}

	// Clear password
	room.Password = ""

	return &room, nil
}

// UpdateRoom mengupdate room
func (s *Service) UpdateRoom(roomID, userID uuid.UUID, req *UpdateRoomRequest) (*models.Room, error) {
	var room models.Room
	if err := s.db.Where("id = ? AND host_id = ?", roomID, userID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to find room for update")
		return nil, fmt.Errorf("internal server error")
	}

	// Validate time range
	if req.StartTime != nil && req.EndTime != nil {
		if req.EndTime.Before(*req.StartTime) {
			return nil, fmt.Errorf("end time must be after start time")
		}
	}

	// Update room fields
	room.Name = req.Name
	room.Description = req.Description
	room.MaxUsers = req.MaxUsers
	room.IsPublic = req.IsPublic
	room.StartTime = req.StartTime
	room.EndTime = req.EndTime

	// Update password if provided
	if req.Password != "" {
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.LogError(err, "Failed to hash room password")
			return nil, fmt.Errorf("failed to process password")
		}
		room.Password = string(hashedBytes)
	}

	if err := s.db.Save(&room).Error; err != nil {
		s.logger.LogError(err, "Failed to update room")
		return nil, fmt.Errorf("failed to update room")
	}

	// Load room with relations
	if err := s.db.Preload("Host").Preload("Settings").First(&room, room.ID).Error; err != nil {
		s.logger.LogError(err, "Failed to load room with relations")
		return nil, fmt.Errorf("failed to load room")
	}

	// Clear password before returning
	room.Password = ""

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("Room updated successfully")
	return &room, nil
}

// DeleteRoom menghapus room
func (s *Service) DeleteRoom(roomID, userID uuid.UUID) error {
	result := s.db.Where("id = ? AND host_id = ?", roomID, userID).Delete(&models.Room{})
	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to delete room")
		return fmt.Errorf("failed to delete room")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("room not found or access denied")
	}

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("Room deleted successfully")
	return nil
}

// JoinRoom bergabung ke room
func (s *Service) JoinRoom(roomID, userID uuid.UUID, req *JoinRoomRequest) (*models.RoomParticipant, error) {
	var room models.Room
	if err := s.db.First(&room, roomID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found")
		}
		s.logger.LogError(err, "Failed to find room for join")
		return nil, fmt.Errorf("internal server error")
	}

	// Check if room is active
	if room.Status != models.RoomStatusActive {
		return nil, fmt.Errorf("room is not active")
	}

	// Check if room has reached max participants
	var participantCount int64
	if err := s.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND status = ?", roomID, models.ParticipantStatusJoined).
		Count(&participantCount).Error; err != nil {
		s.logger.LogError(err, "Failed to count participants")
		return nil, fmt.Errorf("internal server error")
	}

	if int(participantCount) >= room.MaxUsers {
		return nil, fmt.Errorf("room is full")
	}

	// Verify password if room is private
	if room.Password != "" {
		if req.Password == "" {
			return nil, fmt.Errorf("password required")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(room.Password), []byte(req.Password)); err != nil {
			return nil, fmt.Errorf("invalid password")
		}
	}

	// Check if user is already in room
	var existingParticipant models.RoomParticipant
	if err := s.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&existingParticipant).Error; err == nil {
		if existingParticipant.Status == models.ParticipantStatusJoined {
			return nil, fmt.Errorf("already joined room")
		}
		// Rejoin if left before
		existingParticipant.Status = models.ParticipantStatusJoined
		existingParticipant.JoinedAt = time.Now()
		existingParticipant.LeftAt = nil
		if err := s.db.Save(&existingParticipant).Error; err != nil {
			s.logger.LogError(err, "Failed to rejoin room")
			return nil, fmt.Errorf("failed to join room")
		}
		return &existingParticipant, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.LogError(err, "Failed to check existing participant")
		return nil, fmt.Errorf("internal server error")
	}

	// Create new participant
	participant := &models.RoomParticipant{
		RoomID:    roomID,
		UserID:    userID,
		Role:      models.ParticipantRoleParticipant,
		Status:    models.ParticipantStatusJoined,
		JoinedAt:  time.Now(),
		IsVideoOn: true,
	}

	if err := s.db.Create(participant).Error; err != nil {
		s.logger.LogError(err, "Failed to join room")
		return nil, fmt.Errorf("failed to join room")
	}

	// Load participant with user data
	if err := s.db.Preload("User").First(participant, participant.ID).Error; err != nil {
		s.logger.LogError(err, "Failed to load participant with user")
		return nil, fmt.Errorf("failed to load participant")
	}

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("User joined room successfully")
	return participant, nil
}

// LeaveRoom meninggalkan room
func (s *Service) LeaveRoom(roomID, userID uuid.UUID) error {
	var participant models.RoomParticipant
	if err := s.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&participant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not in room")
		}
		s.logger.LogError(err, "Failed to find participant for leave")
		return fmt.Errorf("internal server error")
	}

	// Update participant status
	now := time.Now()
	participant.Status = models.ParticipantStatusLeft
	participant.LeftAt = &now

	if err := s.db.Save(&participant).Error; err != nil {
		s.logger.LogError(err, "Failed to leave room")
		return fmt.Errorf("failed to leave room")
	}

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("User left room successfully")
	return nil
}

// EndRoom mengakhiri room
func (s *Service) EndRoom(roomID, userID uuid.UUID) error {
	var room models.Room
	if err := s.db.Where("id = ? AND host_id = ?", roomID, userID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to find room for end")
		return fmt.Errorf("internal server error")
	}

	// Update room status
	now := time.Now()
	room.Status = models.RoomStatusEnded
	room.EndTime = &now

	if err := s.db.Save(&room).Error; err != nil {
		s.logger.LogError(err, "Failed to end room")
		return fmt.Errorf("failed to end room")
	}

	// Remove all participants from room
	if err := s.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND status = ?", roomID, models.ParticipantStatusJoined).
		Updates(map[string]interface{}{
			"status":  models.ParticipantStatusLeft,
			"left_at": &now,
		}).Error; err != nil {
		s.logger.LogError(err, "Failed to remove participants from room")
	}

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("Room ended successfully")
	return nil
}

// GetRoomParticipants mengambil daftar peserta room
func (s *Service) GetRoomParticipants(roomID uuid.UUID, userID uuid.UUID) ([]*models.RoomParticipant, error) {
	// Check if user has access to room
	var room models.Room
	if err := s.db.Where("id = ?", roomID).
		Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
			true, userID, userID).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to check room access")
		return nil, fmt.Errorf("internal server error")
	}

	var participants []*models.RoomParticipant
	if err := s.db.Where("room_id = ? AND status = ?", roomID, models.ParticipantStatusJoined).
		Preload("User").
		Order("joined_at ASC").
		Find(&participants).Error; err != nil {
		s.logger.LogError(err, "Failed to get room participants")
		return nil, fmt.Errorf("internal server error")
	}

	return participants, nil
}

// GetRoomMessages mengambil pesan dalam room
func (s *Service) GetRoomMessages(roomID uuid.UUID, userID uuid.UUID, page, perPage int) ([]*models.RoomMessage, int64, error) {
	// Check if user has access to room
	var room models.Room
	if err := s.db.Where("id = ?", roomID).
		Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
			true, userID, userID).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to check room access")
		return nil, 0, fmt.Errorf("internal server error")
	}

	var messages []*models.RoomMessage
	var total int64

	query := s.db.Model(&models.RoomMessage{}).Where("room_id = ? AND is_deleted = ?", roomID, false)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		s.logger.LogError(err, "Failed to count room messages")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Get messages with pagination
	offset := (page - 1) * perPage
	if err := query.Preload("Sender").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&messages).Error; err != nil {
		s.logger.LogError(err, "Failed to get room messages")
		return nil, 0, fmt.Errorf("internal server error")
	}

	return messages, total, nil
}

// GetRoomSettings mengambil pengaturan room
func (s *Service) GetRoomSettings(roomID uuid.UUID, userID uuid.UUID) (*models.RoomSetting, error) {
	// Check if user has access to room
	var room models.Room
	if err := s.db.Where("id = ?", roomID).
		Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
			true, userID, userID).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to check room access")
		return nil, fmt.Errorf("internal server error")
	}

	var settings models.RoomSetting
	if err := s.db.Where("room_id = ?", roomID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default settings if not found
			settings = models.RoomSetting{
				RoomID:              roomID,
				AllowScreenShare:    true,
				AllowChat:           true,
				AllowFileShare:      true,
				RequirePassword:     room.Password != "",
				WaitingRoom:         false,
				AutoRecord:          false,
				MaxParticipants:     room.MaxUsers,
				VideoQuality:        "hd",
				AudioQuality:        "high",
				EnableBreakoutRooms: false,
				EnablePolling:       false,
				EnableWhiteboard:    false,
				EnableRecording:     true,
			}
			if err := s.db.Create(&settings).Error; err != nil {
				s.logger.LogError(err, "Failed to create default room settings")
				return nil, fmt.Errorf("failed to create settings")
			}
		} else {
			s.logger.LogError(err, "Failed to get room settings")
			return nil, fmt.Errorf("internal server error")
		}
	}

	return &settings, nil
}

// UpdateRoomSettings mengupdate pengaturan room
func (s *Service) UpdateRoomSettings(roomID, userID uuid.UUID, req *UpdateRoomSettingsRequest) (*models.RoomSetting, error) {
	// Check if user is host
	var room models.Room
	if err := s.db.Where("id = ? AND host_id = ?", roomID, userID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to find room for settings update")
		return nil, fmt.Errorf("internal server error")
	}

	var settings models.RoomSetting
	if err := s.db.Where("room_id = ?", roomID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new settings if not found
			settings = models.RoomSetting{
				RoomID: roomID,
			}
		} else {
			s.logger.LogError(err, "Failed to find room settings")
			return nil, fmt.Errorf("internal server error")
		}
	}

	// Update settings fields
	settings.AllowScreenShare = req.AllowScreenShare
	settings.AllowChat = req.AllowChat
	settings.AllowFileShare = req.AllowFileShare
	settings.RequirePassword = req.RequirePassword
	settings.WaitingRoom = req.WaitingRoom
	settings.AutoRecord = req.AutoRecord
	settings.MaxParticipants = req.MaxParticipants
	settings.VideoQuality = req.VideoQuality
	settings.AudioQuality = req.AudioQuality
	settings.EnableBreakoutRooms = req.EnableBreakoutRooms
	settings.EnablePolling = req.EnablePolling
	settings.EnableWhiteboard = req.EnableWhiteboard
	settings.EnableRecording = req.EnableRecording

	if err := s.db.Save(&settings).Error; err != nil {
		s.logger.LogError(err, "Failed to update room settings")
		return nil, fmt.Errorf("failed to update settings")
	}

	s.logger.WithUserID(userID.String()).WithField("room_id", roomID.String()).Info("Room settings updated successfully")
	return &settings, nil
}

// Helper functions

// isValidRoomType memvalidasi room type
func isValidRoomType(roomType models.RoomType) bool {
	switch roomType {
	case models.RoomTypeMeeting, models.RoomTypeWebinar, models.RoomTypeConference, models.RoomTypeClassroom:
		return true
	default:
		return false
	}
}

// generateRoomCode menggenerate unique room code
func (s *Service) generateRoomCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	for i := 0; i < 10; i++ { // Try 10 times
		code := make([]byte, length)
		for j := 0; j < length; j++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", err
			}
			code[j] = charset[num.Int64()]
		}

		roomCode := string(code)

		// Check if code already exists
		var count int64
		if err := s.db.Model(&models.Room{}).Where("room_code = ?", roomCode).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return roomCode, nil
		}
	}

	return "", errors.New("failed to generate unique room code after 10 attempts")
}

// KickParticipant mengeluarkan peserta dari room (host only)
func (s *Service) KickParticipant(roomID, hostID, participantID uuid.UUID) error {
	// Check if user is host
	var room models.Room
	if err := s.db.Where("id = ? AND host_id = ?", roomID, hostID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to find room for kick participant")
		return fmt.Errorf("internal server error")
	}

	// Don't allow kicking host
	if participantID == hostID {
		return fmt.Errorf("cannot kick host")
	}

	// Update participant status
	now := time.Now()
	result := s.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND user_id = ?", roomID, participantID).
		Updates(map[string]interface{}{
			"status":  models.ParticipantStatusKicked,
			"left_at": &now,
		})

	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to kick participant")
		return fmt.Errorf("failed to kick participant")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("participant not found")
	}

	s.logger.WithUserID(hostID.String()).WithField("room_id", roomID.String()).WithField("participant_id", participantID.String()).Info("Participant kicked successfully")
	return nil
}

// GetRoomStats mengambil statistik room
func (s *Service) GetRoomStats(roomID uuid.UUID, userID uuid.UUID) (map[string]interface{}, error) {
	// Check if user has access to room
	var room models.Room
	if err := s.db.Where("id = ?", roomID).
		Where("is_public = ? OR host_id = ? OR id IN (SELECT room_id FROM room_participants WHERE user_id = ?)",
			true, userID, userID).
		First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("room not found or access denied")
		}
		s.logger.LogError(err, "Failed to check room access")
		return nil, fmt.Errorf("internal server error")
	}

	stats := make(map[string]interface{})

	// Current participant count
	var currentParticipants int64
	if err := s.db.Model(&models.RoomParticipant{}).
		Where("room_id = ? AND status = ?", roomID, models.ParticipantStatusJoined).
		Count(&currentParticipants).Error; err != nil {
		s.logger.LogError(err, "Failed to count current participants")
		return nil, fmt.Errorf("internal server error")
	}
	stats["current_participants"] = currentParticipants

	// Total participant count (including left)
	var totalParticipants int64
	if err := s.db.Model(&models.RoomParticipant{}).
		Where("room_id = ?", roomID).
		Count(&totalParticipants).Error; err != nil {
		s.logger.LogError(err, "Failed to count total participants")
		return nil, fmt.Errorf("internal server error")
	}
	stats["total_participants"] = totalParticipants

	// Message count
	var messageCount int64
	if err := s.db.Model(&models.RoomMessage{}).
		Where("room_id = ? AND is_deleted = ?", roomID, false).
		Count(&messageCount).Error; err != nil {
		s.logger.LogError(err, "Failed to count messages")
		return nil, fmt.Errorf("internal server error")
	}
	stats["message_count"] = messageCount

	// Room duration
	if room.StartTime != nil {
		if room.EndTime != nil {
			duration := room.EndTime.Sub(*room.StartTime)
			stats["duration_minutes"] = int(duration.Minutes())
		} else {
			duration := time.Since(*room.StartTime)
			stats["duration_minutes"] = int(duration.Minutes())
		}
	} else {
		stats["duration_minutes"] = 0
	}

	stats["room_status"] = room.Status
	stats["room_type"] = room.Type

	return stats, nil
}
