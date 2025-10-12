package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Room model untuk tabel rooms
type Room struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	HostID      uuid.UUID      `json:"host_id" gorm:"type:uuid;not null"`
	RoomCode    string         `json:"room_code" gorm:"uniqueIndex;not null"`
	Password    string         `json:"-"` // password untuk room private
	MaxUsers    int            `json:"max_users" gorm:"default:50"`
	Status      RoomStatus     `json:"status" gorm:"default:'active'"`
	Type        RoomType       `json:"type" gorm:"default:'meeting'"`
	IsPublic    bool           `json:"is_public" gorm:"default:true"`
	IsRecording bool           `json:"is_recording" gorm:"default:false"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Host         *User             `json:"host,omitempty" gorm:"foreignKey:HostID"`
	Participants []RoomParticipant `json:"participants,omitempty"`
	Messages     []RoomMessage     `json:"messages,omitempty"`
	Settings     *RoomSetting      `json:"settings,omitempty" gorm:"foreignKey:RoomID"`
	Histories    []MeetingHistory  `json:"histories,omitempty" gorm:"foreignKey:RoomID"`
}

// RoomStatus enum untuk status room
type RoomStatus string

const (
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
	RoomStatusEnded    RoomStatus = "ended"
	RoomStatusBlocked  RoomStatus = "blocked"
)

// RoomType enum untuk tipe room
type RoomType string

const (
	RoomTypeMeeting    RoomType = "meeting"
	RoomTypeWebinar    RoomType = "webinar"
	RoomTypeConference RoomType = "conference"
	RoomTypeClassroom  RoomType = "classroom"
)

// RoomParticipant model untuk tabel room_participants
type RoomParticipant struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID          uuid.UUID         `json:"room_id" gorm:"type:uuid;not null"`
	UserID          uuid.UUID         `json:"user_id" gorm:"type:uuid;not null"`
	Role            ParticipantRole   `json:"role" gorm:"default:'participant'"`
	Status          ParticipantStatus `json:"status" gorm:"default:'joined'"`
	JoinedAt        time.Time         `json:"joined_at"`
	LeftAt          *time.Time        `json:"left_at"`
	IsMuted         bool              `json:"is_muted" gorm:"default:false"`
	IsVideoOn       bool              `json:"is_video_on" gorm:"default:true"`
	IsScreenSharing bool              `json:"is_screen_sharing" gorm:"default:false"`
	HandRaised      bool              `json:"hand_raised" gorm:"default:false"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	// Relations
	Room *Room `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ParticipantRole enum untuk role participant
type ParticipantRole string

const (
	ParticipantRoleHost        ParticipantRole = "host"
	ParticipantRoleModerator   ParticipantRole = "moderator"
	ParticipantRoleParticipant ParticipantRole = "participant"
)

// ParticipantStatus enum untuk status participant
type ParticipantStatus string

const (
	ParticipantStatusJoined ParticipantStatus = "joined"
	ParticipantStatusLeft   ParticipantStatus = "left"
	ParticipantStatusKicked ParticipantStatus = "kicked"
	ParticipantStatusBanned ParticipantStatus = "banned"
)

// RoomMessage model untuk tabel room_messages
type RoomMessage struct {
	ID        uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID    uuid.UUID   `json:"room_id" gorm:"type:uuid;not null"`
	SenderID  uuid.UUID   `json:"sender_id" gorm:"type:uuid;not null"`
	Message   string      `json:"message" gorm:"not null"`
	Type      MessageType `json:"type" gorm:"default:'text'"`
	FileURL   string      `json:"file_url"`
	FileName  string      `json:"file_name"`
	FileSize  int64       `json:"file_size"`
	IsDeleted bool        `json:"is_deleted" gorm:"default:false"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	// Relations
	Room   *Room `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Sender *User `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
}

// MessageType enum untuk tipe message
type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeFile   MessageType = "file"
	MessageTypeImage  MessageType = "image"
	MessageTypeSystem MessageType = "system"
)

// RoomSetting model untuk tabel room_settings
type RoomSetting struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID              uuid.UUID `json:"room_id" gorm:"type:uuid;not null;uniqueIndex"`
	AllowScreenShare    bool      `json:"allow_screen_share" gorm:"default:true"`
	AllowChat           bool      `json:"allow_chat" gorm:"default:true"`
	AllowFileShare      bool      `json:"allow_file_share" gorm:"default:true"`
	RequirePassword     bool      `json:"require_password" gorm:"default:false"`
	WaitingRoom         bool      `json:"waiting_room" gorm:"default:false"`
	AutoRecord          bool      `json:"auto_record" gorm:"default:false"`
	MaxParticipants     int       `json:"max_participants" gorm:"default:50"`
	VideoQuality        string    `json:"video_quality" gorm:"default:'hd'"`
	AudioQuality        string    `json:"audio_quality" gorm:"default:'high'"`
	EnableBreakoutRooms bool      `json:"enable_breakout_rooms" gorm:"default:false"`
	EnablePolling       bool      `json:"enable_polling" gorm:"default:false"`
	EnableWhiteboard    bool      `json:"enable_whiteboard" gorm:"default:false"`
	EnableRecording     bool      `json:"enable_recording" gorm:"default:true"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relations
	Room *Room `json:"room,omitempty" gorm:"foreignKey:RoomID"`
}

// MeetingHistory model untuk tabel meeting_history
type MeetingHistory struct {
	ID               uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RoomID           uuid.UUID     `json:"room_id" gorm:"type:uuid;not null"`
	HostID           uuid.UUID     `json:"host_id" gorm:"type:uuid;not null"`
	Title            string        `json:"title" gorm:"not null"`
	Description      string        `json:"description"`
	StartTime        time.Time     `json:"start_time" gorm:"not null"`
	EndTime          *time.Time    `json:"end_time"`
	Duration         int           `json:"duration"` // in minutes
	ParticipantCount int           `json:"participant_count" gorm:"default:0"`
	RecordingURL     string        `json:"recording_url"`
	RecordingSize    int64         `json:"recording_size"`
	Status           MeetingStatus `json:"status" gorm:"default:'scheduled'"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`

	// Relations
	Room *Room `json:"room,omitempty" gorm:"foreignKey:RoomID"`
	Host *User `json:"host,omitempty" gorm:"foreignKey:HostID"`
}

// MeetingStatus enum untuk status meeting
type MeetingStatus string

const (
	MeetingStatusScheduled MeetingStatus = "scheduled"
	MeetingStatusOngoing   MeetingStatus = "ongoing"
	MeetingStatusEnded     MeetingStatus = "ended"
	MeetingStatusCancelled MeetingStatus = "cancelled"
)

// Notification model untuk tabel notifications
type Notification struct {
	ID        uuid.UUID        `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:uuid;not null"`
	Title     string           `json:"title" gorm:"not null"`
	Message   string           `json:"message" gorm:"not null"`
	Type      NotificationType `json:"type" gorm:"not null"`
	Data      string           `json:"data"` // JSON data for additional info
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// NotificationType enum untuk tipe notification
type NotificationType string

const (
	NotificationTypeRoomInvite     NotificationType = "room_invite"
	NotificationTypeRoomMessage    NotificationType = "room_message"
	NotificationTypeMeetingStart   NotificationType = "meeting_start"
	NotificationTypeMeetingEnd     NotificationType = "meeting_end"
	NotificationTypeSystem         NotificationType = "system"
	NotificationTypeContactRequest NotificationType = "contact_request"
)

// TableName methods
func (Room) TableName() string {
	return "rooms"
}

func (RoomParticipant) TableName() string {
	return "room_participants"
}

func (RoomMessage) TableName() string {
	return "room_messages"
}

func (RoomSetting) TableName() string {
	return "room_settings"
}

func (MeetingHistory) TableName() string {
	return "meeting_history"
}

func (Notification) TableName() string {
	return "notifications"
}

// BeforeCreate hooks
func (r *Room) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (rp *RoomParticipant) BeforeCreate(tx *gorm.DB) error {
	if rp.ID == uuid.Nil {
		rp.ID = uuid.New()
	}
	return nil
}

func (rm *RoomMessage) BeforeCreate(tx *gorm.DB) error {
	if rm.ID == uuid.Nil {
		rm.ID = uuid.New()
	}
	return nil
}

func (rs *RoomSetting) BeforeCreate(tx *gorm.DB) error {
	if rs.ID == uuid.Nil {
		rs.ID = uuid.New()
	}
	return nil
}

func (mh *MeetingHistory) BeforeCreate(tx *gorm.DB) error {
	if mh.ID == uuid.Nil {
		mh.ID = uuid.New()
	}
	return nil
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}

// Helper methods
func (r *Room) IsActive() bool {
	return r.Status == RoomStatusActive
}

func (r *Room) IsOngoing() bool {
	if r.StartTime == nil {
		return false
	}
	now := time.Now()
	if r.EndTime != nil {
		return now.After(*r.StartTime) && now.Before(*r.EndTime)
	}
	return now.After(*r.StartTime)
}

func (rp *RoomParticipant) IsHost() bool {
	return rp.Role == ParticipantRoleHost
}

func (rp *RoomParticipant) IsModerator() bool {
	return rp.Role == ParticipantRoleModerator || rp.Role == ParticipantRoleHost
}

func (rp *RoomParticipant) IsActive() bool {
	return rp.Status == ParticipantStatusJoined
}

func (mh *MeetingHistory) IsOngoing() bool {
	return mh.Status == MeetingStatusOngoing
}

func (mh *MeetingHistory) GetDurationMinutes() int {
	if mh.EndTime == nil {
		return int(time.Since(mh.StartTime).Minutes())
	}
	return int(mh.EndTime.Sub(mh.StartTime).Minutes())
}
