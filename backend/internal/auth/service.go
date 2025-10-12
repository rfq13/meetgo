package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/webrtc-meeting/backend/internal/config"
	"github.com/webrtc-meeting/backend/models"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// Service struct untuk authentication service
type Service struct {
	db     *gorm.DB
	config *config.Config
	logger *logger.Logger
}

// NewService membuat authentication service baru
func NewService(db *gorm.DB, cfg *config.Config, log *logger.Logger) *Service {
	return &Service{
		db:     db,
		config: cfg,
		logger: log,
	}
}

// LoginRequest struct untuk request login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest struct untuk request register
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,min=1,max=50"`
	LastName  string `json:"last_name" binding:"required,min=1,max=50"`
	Phone     string `json:"phone"`
}

// LoginResponse struct untuk response login
type LoginResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

// RefreshTokenRequest struct untuk request refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest struct untuk request change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Claims struct untuk JWT claims
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// Login melakukan authentication user
func (s *Service) Login(req *LoginRequest, clientIP string) (*LoginResponse, error) {
	// Cari user berdasarkan email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.LogAuthEvent("login_failed", "", clientIP, false)
			return nil, fmt.Errorf("invalid credentials")
		}
		s.logger.LogError(err, "Failed to find user during login")
		return nil, fmt.Errorf("internal server error")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.LogAuthEvent("login_failed", user.ID.String(), clientIP, false)
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		s.logger.LogAuthEvent("login_blocked", user.ID.String(), clientIP, false)
		return nil, fmt.Errorf("user account is not active")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		s.logger.LogError(err, "Failed to generate access token")
		return nil, fmt.Errorf("failed to generate token")
	}

	refreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		s.logger.LogError(err, "Failed to generate refresh token")
		return nil, fmt.Errorf("failed to generate token")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.db.Save(&user).Error; err != nil {
		s.logger.LogError(err, "Failed to update last login")
	}

	// Create session
	session := &models.UserSession{
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IPAddress:    clientIP,
		ExpiresAt:    time.Now().Add(s.config.JWT.ExpirationTime),
	}
	if err := s.db.Create(session).Error; err != nil {
		s.logger.LogError(err, "Failed to create user session")
	}

	s.logger.LogAuthEvent("login_success", user.ID.String(), clientIP, true)

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWT.ExpirationTime.Seconds()),
	}, nil
}

// Register melakukan registrasi user baru
func (s *Service) Register(req *RegisterRequest, clientIP string) (*models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.LogError(err, "Failed to check existing email")
		return nil, fmt.Errorf("internal server error")
	}

	// Check if username already exists
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.LogError(err, "Failed to check existing username")
		return nil, fmt.Errorf("internal server error")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.LogError(err, "Failed to hash password")
		return nil, fmt.Errorf("internal server error")
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Status:    models.UserStatusActive,
		Role:      models.UserRoleUser,
	}

	if err := s.db.Create(user).Error; err != nil {
		s.logger.LogError(err, "Failed to create user")
		return nil, fmt.Errorf("failed to create user")
	}

	// Create user settings
	userSettings := &models.UserSetting{
		UserID: user.ID,
	}
	if err := s.db.Create(userSettings).Error; err != nil {
		s.logger.LogError(err, "Failed to create user settings")
	}

	s.logger.LogAuthEvent("register_success", user.ID.String(), clientIP, true)

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// RefreshToken melakukan refresh access token
func (s *Service) RefreshToken(req *RefreshTokenRequest, clientIP string) (*LoginResponse, error) {
	// Validate refresh token
	claims, err := s.validateToken(req.RefreshToken)
	if err != nil {
		s.logger.LogAuthEvent("refresh_token_invalid", "", clientIP, false)
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Find session
	var session models.UserSession
	if err := s.db.Where("refresh_token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.LogAuthEvent("refresh_token_expired", claims.UserID.String(), clientIP, false)
			return nil, fmt.Errorf("refresh token expired")
		}
		s.logger.LogError(err, "Failed to find session")
		return nil, fmt.Errorf("internal server error")
	}

	// Find user
	var user models.User
	if err := s.db.First(&user, session.UserID).Error; err != nil {
		s.logger.LogError(err, "Failed to find user during refresh")
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !user.IsActive() {
		s.logger.LogAuthEvent("refresh_token_blocked", user.ID.String(), clientIP, false)
		return nil, fmt.Errorf("user account is not active")
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(&user)
	if err != nil {
		s.logger.LogError(err, "Failed to generate new access token")
		return nil, fmt.Errorf("failed to generate token")
	}

	newRefreshToken, err := s.generateRefreshToken(&user)
	if err != nil {
		s.logger.LogError(err, "Failed to generate new refresh token")
		return nil, fmt.Errorf("failed to generate token")
	}

	// Update session
	session.Token = accessToken
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(s.config.JWT.ExpirationTime)
	if err := s.db.Save(&session).Error; err != nil {
		s.logger.LogError(err, "Failed to update session")
	}

	s.logger.LogAuthEvent("refresh_token_success", user.ID.String(), clientIP, true)

	return &LoginResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.config.JWT.ExpirationTime.Seconds()),
	}, nil
}

// Logout melakukan logout user
func (s *Service) Logout(token string, clientIP string) error {
	// Find and delete session
	if err := s.db.Where("token = ?", token).Delete(&models.UserSession{}).Error; err != nil {
		s.logger.LogError(err, "Failed to delete session during logout")
		return fmt.Errorf("failed to logout")
	}

	s.logger.LogAuthEvent("logout_success", "", clientIP, true)
	return nil
}

// ChangePassword mengubah password user
func (s *Service) ChangePassword(userID uuid.UUID, req *ChangePasswordRequest) error {
	// Find user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		s.logger.LogError(err, "Failed to find user during password change")
		return fmt.Errorf("user not found")
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

	// Delete all sessions (force re-login)
	if err := s.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error; err != nil {
		s.logger.LogError(err, "Failed to delete user sessions after password change")
	}

	s.logger.WithUserID(userID.String()).Info("Password changed successfully")
	return nil
}

// ValidateToken validates access token and returns user info
func (s *Service) ValidateToken(tokenString string) (*models.User, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Find user
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, fmt.Errorf("user account is not active")
	}

	// Check if session exists and is valid
	var session models.UserSession
	if err := s.db.Where("token = ? AND expires_at > ?", tokenString, time.Now()).First(&session).Error; err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	return &user, nil
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

// generateAccessToken generates JWT access token
func (s *Service) generateAccessToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "webrtc-meeting",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// generateRefreshToken generates JWT refresh token
func (s *Service) generateRefreshToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "webrtc-meeting",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.RefreshSecret))
}

// validateToken validates JWT token and returns claims
func (s *Service) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// CleanupExpiredSessions cleans up expired sessions
func (s *Service) CleanupExpiredSessions() error {
	result := s.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{})
	if result.Error != nil {
		s.logger.LogError(result.Error, "Failed to cleanup expired sessions")
		return result.Error
	}

	if result.RowsAffected > 0 {
		s.logger.Infof("Cleaned up %d expired sessions", result.RowsAffected)
	}

	return nil
}
