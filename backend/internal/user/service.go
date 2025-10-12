package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/webrtc-meeting/backend/models"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Service struct untuk user service
type Service struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewService membuat user service baru
func NewService(db *gorm.DB, log *logger.Logger) *Service {
	return &Service{
		db:     db,
		logger: log,
	}
}

// UpdateProfileRequest struct untuk request update profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
}

// ChangePasswordRequest struct untuk request change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// AddContactRequest struct untuk request add contact
type AddContactRequest struct {
	ContactID string `json:"contact_id" binding:"required"`
}

// UpdateSettingsRequest struct untuk request update settings
type UpdateSettingsRequest struct {
	Language             string `json:"language"`
	Timezone             string `json:"timezone"`
	EmailNotifications   bool   `json:"email_notifications"`
	PushNotifications    bool   `json:"push_notifications"`
	VideoQuality         string `json:"video_quality"`
	AudioQuality         string `json:"audio_quality"`
	AutoJoinMicrophone   bool   `json:"auto_join_microphone"`
	AutoJoinCamera       bool   `json:"auto_join_camera"`
	ScreenShareQuality   string `json:"screen_share_quality"`
	RecordingEnabled     bool   `json:"recording_enabled"`
	ChatEnabled          bool   `json:"chat_enabled"`
	ParticipantListShown bool   `json:"participant_list_shown"`
}

// GetProfile mengambil profile user
func (s *Service) GetProfile(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Settings").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		s.logger.LogError(err, "Failed to get user profile")
		return nil, fmt.Errorf("internal server error")
	}

	// Clear password
	user.Password = ""
	return &user, nil
}

// UpdateProfile mengupdate profile user
func (s *Service) UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		s.logger.LogError(err, "Failed to find user for profile update")
		return nil, fmt.Errorf("internal server error")
	}

	// Update user fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Avatar = req.Avatar
	user.Phone = req.Phone

	if err := s.db.Save(&user).Error; err != nil {
		s.logger.LogError(err, "Failed to update user profile")
		return nil, fmt.Errorf("failed to update profile")
	}

	s.logger.WithUserID(userID.String()).Info("User profile updated successfully")

	// Clear password before returning
	user.Password = ""
	return &user, nil
}

// ChangePassword mengubah password user
func (s *Service) ChangePassword(userID uuid.UUID, req *ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		s.logger.LogError(err, "Failed to find user for password change")
		return fmt.Errorf("internal server error")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.LogError(err, "Failed to hash new password")
		return fmt.Errorf("internal server error")
	}

	// Update password
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		s.logger.LogError(err, "Failed to update password")
		return fmt.Errorf("failed to update password")
	}

	s.logger.WithUserID(userID.String()).Info("User password changed successfully")
	return nil
}

// GetContacts mengambil daftar kontak user
func (s *Service) GetContacts(userID uuid.UUID, page, perPage int, search string) ([]*models.UserContact, int64, error) {
	var contacts []*models.UserContact
	var total int64

	query := s.db.Where("user_id = ?", userID)

	// Add search filter
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Joins("JOIN users ON user_contacts.contact_id = users.id").
			Where("(LOWER(users.first_name) LIKE ? OR LOWER(users.last_name) LIKE ? OR LOWER(users.email) LIKE ?)", searchTerm, searchTerm, searchTerm)
	}

	// Count total
	if err := query.Model(&models.UserContact{}).Count(&total).Error; err != nil {
		s.logger.LogError(err, "Failed to count user contacts")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Get contacts with pagination
	offset := (page - 1) * perPage
	if err := query.Preload("Contact").Offset(offset).Limit(perPage).Find(&contacts).Error; err != nil {
		s.logger.LogError(err, "Failed to get user contacts")
		return nil, 0, fmt.Errorf("internal server error")
	}

	return contacts, total, nil
}

// AddContact menambahkan kontak baru
func (s *Service) AddContact(userID, contactID uuid.UUID) (*models.UserContact, error) {
	// Check if contact exists
	var contact models.User
	if err := s.db.First(&contact, contactID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contact user not found")
		}
		s.logger.LogError(err, "Failed to find contact user")
		return nil, fmt.Errorf("internal server error")
	}

	// Check if contact already exists
	var existingContact models.UserContact
	if err := s.db.Where("user_id = ? AND contact_id = ?", userID, contactID).First(&existingContact).Error; err == nil {
		return nil, fmt.Errorf("contact already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.LogError(err, "Failed to check existing contact")
		return nil, fmt.Errorf("internal server error")
	}

	// Create contact
	userContact := &models.UserContact{
		UserID:    userID,
		ContactID: contactID,
		Status:    "accepted",
	}

	if err := s.db.Create(userContact).Error; err != nil {
		s.logger.LogError(err, "Failed to add contact")
		return nil, fmt.Errorf("failed to add contact")
	}

	// Preload contact data
	if err := s.db.Preload("Contact").First(userContact, userContact.ID).Error; err != nil {
		s.logger.LogError(err, "Failed to preload contact data")
		return nil, fmt.Errorf("failed to get contact data")
	}

	s.logger.WithUserID(userID.String()).WithField("contact_id", contactID.String()).Info("Contact added successfully")
	return userContact, nil
}

// RemoveContact menghapus kontak
func (s *Service) RemoveContact(userID, contactID uuid.UUID) error {
	result := s.db.Where("user_id = ? AND contact_id = ?", userID, contactID).Delete(&models.UserContact{})
	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to remove contact")
		return fmt.Errorf("failed to remove contact")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("contact not found")
	}

	s.logger.WithUserID(userID.String()).WithField("contact_id", contactID.String()).Info("Contact removed successfully")
	return nil
}

// SearchUsers mencari user untuk ditambahkan sebagai kontak
func (s *Service) SearchUsers(userID uuid.UUID, query string, page, perPage int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	searchTerm := "%" + strings.ToLower(query) + "%"

	// Build query to search users excluding current user and existing contacts
	baseQuery := s.db.Where("LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(email) LIKE ?", searchTerm, searchTerm, searchTerm).
		Where("id != ?", userID).
		Where("status = ?", models.UserStatusActive)

	// Count total
	if err := baseQuery.Model(&models.User{}).Count(&total).Error; err != nil {
		s.logger.LogError(err, "Failed to count search results")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Get users with pagination
	offset := (page - 1) * perPage
	if err := baseQuery.Offset(offset).Limit(perPage).Find(&users).Error; err != nil {
		s.logger.LogError(err, "Failed to search users")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Clear passwords
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}

// GetSettings mengambil pengaturan user
func (s *Service) GetSettings(userID uuid.UUID) (*models.UserSetting, error) {
	var settings models.UserSetting
	if err := s.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create default settings if not found
			settings = models.UserSetting{
				UserID:               userID,
				Language:             "en",
				Timezone:             "UTC",
				EmailNotifications:   true,
				PushNotifications:    true,
				VideoQuality:         "hd",
				AudioQuality:         "high",
				AutoJoinMicrophone:   false,
				AutoJoinCamera:       false,
				ScreenShareQuality:   "hd",
				RecordingEnabled:     true,
				ChatEnabled:          true,
				ParticipantListShown: true,
			}
			if err := s.db.Create(&settings).Error; err != nil {
				s.logger.LogError(err, "Failed to create default user settings")
				return nil, fmt.Errorf("failed to create settings")
			}
		} else {
			s.logger.LogError(err, "Failed to get user settings")
			return nil, fmt.Errorf("internal server error")
		}
	}

	return &settings, nil
}

// UpdateSettings mengupdate pengaturan user
func (s *Service) UpdateSettings(userID uuid.UUID, req *UpdateSettingsRequest) (*models.UserSetting, error) {
	var settings models.UserSetting
	if err := s.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new settings if not found
			settings = models.UserSetting{
				UserID: userID,
			}
		} else {
			s.logger.LogError(err, "Failed to find user settings")
			return nil, fmt.Errorf("internal server error")
		}
	}

	// Update settings fields
	if req.Language != "" {
		settings.Language = req.Language
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}
	settings.EmailNotifications = req.EmailNotifications
	settings.PushNotifications = req.PushNotifications
	if req.VideoQuality != "" {
		settings.VideoQuality = req.VideoQuality
	}
	if req.AudioQuality != "" {
		settings.AudioQuality = req.AudioQuality
	}
	settings.AutoJoinMicrophone = req.AutoJoinMicrophone
	settings.AutoJoinCamera = req.AutoJoinCamera
	if req.ScreenShareQuality != "" {
		settings.ScreenShareQuality = req.ScreenShareQuality
	}
	settings.RecordingEnabled = req.RecordingEnabled
	settings.ChatEnabled = req.ChatEnabled
	settings.ParticipantListShown = req.ParticipantListShown

	if err := s.db.Save(&settings).Error; err != nil {
		s.logger.LogError(err, "Failed to update user settings")
		return nil, fmt.Errorf("failed to update settings")
	}

	s.logger.WithUserID(userID.String()).Info("User settings updated successfully")
	return &settings, nil
}

// GetNotifications mengambil notifikasi user
func (s *Service) GetNotifications(userID uuid.UUID, page, perPage int, isRead *bool) ([]*models.Notification, int64, error) {
	var notifications []*models.Notification
	var total int64

	query := s.db.Where("user_id = ?", userID)

	// Add read filter if specified
	if isRead != nil {
		query = query.Where("is_read = ?", *isRead)
	}

	// Count total
	if err := query.Model(&models.Notification{}).Count(&total).Error; err != nil {
		s.logger.LogError(err, "Failed to count notifications")
		return nil, 0, fmt.Errorf("internal server error")
	}

	// Get notifications with pagination
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&notifications).Error; err != nil {
		s.logger.LogError(err, "Failed to get notifications")
		return nil, 0, fmt.Errorf("internal server error")
	}

	return notifications, total, nil
}

// MarkNotificationAsRead menandai notifikasi sebagai telah dibaca
func (s *Service) MarkNotificationAsRead(userID, notificationID uuid.UUID) error {
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to mark notification as read")
		return fmt.Errorf("failed to mark notification as read")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	s.logger.WithUserID(userID.String()).WithField("notification_id", notificationID.String()).Info("Notification marked as read")
	return nil
}

// MarkAllNotificationsAsRead menandai semua notifikasi sebagai telah dibaca
func (s *Service) MarkAllNotificationsAsRead(userID uuid.UUID) error {
	result := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to mark all notifications as read")
		return fmt.Errorf("failed to mark all notifications as read")
	}

	s.logger.WithUserID(userID.String()).Infof("Marked %d notifications as read", result.RowsAffected)
	return nil
}

// DeleteNotification menghapus notifikasi
func (s *Service) DeleteNotification(userID, notificationID uuid.UUID) error {
	result := s.db.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{})
	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to delete notification")
		return fmt.Errorf("failed to delete notification")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	s.logger.WithUserID(userID.String()).WithField("notification_id", notificationID.String()).Info("Notification deleted successfully")
	return nil
}

// GetUserStats mengambil statistik user
func (s *Service) GetUserStats(userID uuid.UUID) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count rooms created
	var roomsCreated int64
	if err := s.db.Model(&models.Room{}).Where("host_id = ?", userID).Count(&roomsCreated).Error; err != nil {
		s.logger.LogError(err, "Failed to count rooms created")
		return nil, fmt.Errorf("internal server error")
	}
	stats["rooms_created"] = roomsCreated

	// Count meetings participated
	var meetingsParticipated int64
	if err := s.db.Model(&models.RoomParticipant{}).Where("user_id = ?", userID).Count(&meetingsParticipated).Error; err != nil {
		s.logger.LogError(err, "Failed to count meetings participated")
		return nil, fmt.Errorf("internal server error")
	}
	stats["meetings_participated"] = meetingsParticipated

	// Count contacts
	var contactsCount int64
	if err := s.db.Model(&models.UserContact{}).Where("user_id = ?", userID).Count(&contactsCount).Error; err != nil {
		s.logger.LogError(err, "Failed to count contacts")
		return nil, fmt.Errorf("internal server error")
	}
	stats["contacts_count"] = contactsCount

	// Count unread notifications
	var unreadNotifications int64
	if err := s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadNotifications).Error; err != nil {
		s.logger.LogError(err, "Failed to count unread notifications")
		return nil, fmt.Errorf("internal server error")
	}
	stats["unread_notifications"] = unreadNotifications

	// Get total meeting duration
	var totalDuration int64
	if err := s.db.Model(&models.RoomParticipant{}).
		Where("user_id = ? AND left_at IS NOT NULL", userID).
		Select("COALESCE(SUM(EXTRACT(EPOCH FROM (left_at - joined_at))/60), 0)").
		Scan(&totalDuration).Error; err != nil {
		s.logger.LogError(err, "Failed to calculate total meeting duration")
		// Don't return error, just set to 0
		totalDuration = 0
	}
	stats["total_meeting_duration_minutes"] = totalDuration

	return stats, nil
}

// DeactivateAccount menonaktifkan akun user
func (s *Service) DeactivateAccount(userID uuid.UUID) error {
	result := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("status", models.UserStatusInactive)

	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to deactivate user account")
		return fmt.Errorf("failed to deactivate account")
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Delete all sessions
	if err := s.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error; err != nil {
		s.logger.LogError(err, "Failed to delete user sessions after deactivation")
	}

	s.logger.WithUserID(userID.String()).Warn("User account deactivated")
	return nil
}

// UpdateLastLogin update last login time
func (s *Service) UpdateLastLogin(userID uuid.UUID) error {
	now := time.Now()
	result := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login", &now)

	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to update last login")
		return fmt.Errorf("failed to update last login")
	}

	return nil
}
