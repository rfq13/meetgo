package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User model untuk tabel users
type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Avatar    string         `json:"avatar"`
	Phone     string         `json:"phone"`
	Status    UserStatus     `json:"status" gorm:"default:'active'"`
	Role      UserRole       `json:"role" gorm:"default:'user'"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Rooms         []Room         `json:"rooms,omitempty" gorm:"many2many:room_participants;"`
	Sessions      []UserSession  `json:"sessions,omitempty"`
	Contacts      []UserContact  `json:"contacts,omitempty" gorm:"foreignKey:UserID"`
	Settings      *UserSetting   `json:"settings,omitempty" gorm:"foreignKey:UserID"`
	Messages      []RoomMessage  `json:"messages,omitempty" gorm:"foreignKey:SenderID"`
	Notifications []Notification `json:"notifications,omitempty" gorm:"foreignKey:UserID"`
}

// UserStatus enum untuk status user
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBlocked  UserStatus = "blocked"
	UserStatusPending  UserStatus = "pending"
)

// UserRole enum untuk role user
type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
	UserRoleUser      UserRole = "user"
)

// UserSession model untuk tabel user_sessions
type UserSession struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token        string    `json:"token" gorm:"uniqueIndex;not null"`
	RefreshToken string    `json:"refresh_token" gorm:"uniqueIndex;not null"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserContact model untuk tabel user_contacts
type UserContact struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ContactID uuid.UUID `json:"contact_id" gorm:"type:uuid;not null"`
	Nickname  string    `json:"nickname"`
	Status    string    `json:"status" gorm:"default:'pending'"` // pending, accepted, blocked
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	User    *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Contact *User `json:"contact,omitempty" gorm:"foreignKey:ContactID"`
}

// UserSetting model untuk tabel user_settings
type UserSetting struct {
	ID                   uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID               uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Language             string    `json:"language" gorm:"default:'en'"`
	Timezone             string    `json:"timezone" gorm:"default:'UTC'"`
	EmailNotifications   bool      `json:"email_notifications" gorm:"default:true"`
	PushNotifications    bool      `json:"push_notifications" gorm:"default:true"`
	VideoQuality         string    `json:"video_quality" gorm:"default:'hd'"`
	AudioQuality         string    `json:"audio_quality" gorm:"default:'high'"`
	AutoJoinMicrophone   bool      `json:"auto_join_microphone" gorm:"default:false"`
	AutoJoinCamera       bool      `json:"auto_join_camera" gorm:"default:false"`
	ScreenShareQuality   string    `json:"screen_share_quality" gorm:"default:'hd'"`
	RecordingEnabled     bool      `json:"recording_enabled" gorm:"default:true"`
	ChatEnabled          bool      `json:"chat_enabled" gorm:"default:true"`
	ParticipantListShown bool      `json:"participant_list_shown" gorm:"default:true"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName untuk User model
func (User) TableName() string {
	return "users"
}

// TableName untuk UserSession model
func (UserSession) TableName() string {
	return "user_sessions"
}

// TableName untuk UserContact model
func (UserContact) TableName() string {
	return "user_contacts"
}

// TableName untuk UserSetting model
func (UserSetting) TableName() string {
	return "user_settings"
}

// BeforeCreate hook untuk User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook untuk UserSession
func (us *UserSession) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook untuk UserContact
func (uc *UserContact) BeforeCreate(tx *gorm.DB) error {
	if uc.ID == uuid.Nil {
		uc.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook untuk UserSetting
func (us *UserSetting) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

// GetFullName mengembalikan nama lengkap user
func (u *User) GetFullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// IsActive mengembalikan true jika user aktif
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// IsAdmin mengembalikan true jika user adalah admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsModerator mengembalikan true jika user adalah moderator atau admin
func (u *User) IsModerator() bool {
	return u.Role == UserRoleModerator || u.Role == UserRoleAdmin
}

// IsExpired mengembalikan true jika session sudah expired
func (us *UserSession) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}
