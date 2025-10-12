package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/webrtc-meeting/backend/internal/config"
	"github.com/webrtc-meeting/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database struct untuk menyimpan koneksi database
type Database struct {
	DB *gorm.DB
}

// NewDatabase membuat koneksi database baru
func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := cfg.GetDatabaseURL()

	// Konfigurasi GORM logger
	gormLogger := logger.New(
		logrus.New(),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  getLogLevel(cfg.Logger.Level),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Buka koneksi database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)                 // Maksimum koneksi terbuka
	sqlDB.SetMaxIdleConns(25)                 // Maksimum koneksi idle
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maksimum lifetime koneksi
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Maksimum idle time koneksi

	// Test koneksi
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Database connection established successfully")

	database := &Database{
		DB: db,
	}

	// Auto migrate tables
	if err := database.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return database, nil
}

// AutoMigrate melakukan auto migration untuk semua model
func (d *Database) AutoMigrate() error {
	logrus.Info("Starting database migration...")

	// List semua model yang akan di-migrate
	models := []interface{}{
		&models.User{},
		&models.UserSession{},
		&models.UserContact{},
		&models.UserSetting{},
		&models.Room{},
		&models.RoomParticipant{},
		&models.RoomMessage{},
		&models.RoomSetting{},
		&models.MeetingHistory{},
		&models.Notification{},
	}

	// Lakukan migration
	for _, model := range models {
		if err := d.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
		logrus.Infof("Migrated %T successfully", model)
	}

	// Create indexes setelah migration
	if err := d.createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	logrus.Info("Database migration completed successfully")
	return nil
}

// createIndexes membuat indexes yang diperlukan
func (d *Database) createIndexes() error {
	// Index untuk users table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return fmt.Errorf("failed to create idx_users_email: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return fmt.Errorf("failed to create idx_users_username: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)").Error; err != nil {
		return fmt.Errorf("failed to create idx_users_status: %w", err)
	}

	// Index untuk rooms table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_rooms_host_id ON rooms(host_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_rooms_host_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_rooms_room_code ON rooms(room_code)").Error; err != nil {
		return fmt.Errorf("failed to create idx_rooms_room_code: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status)").Error; err != nil {
		return fmt.Errorf("failed to create idx_rooms_status: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_rooms_type ON rooms(type)").Error; err != nil {
		return fmt.Errorf("failed to create idx_rooms_type: %w", err)
	}

	// Index untuk room_participants table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_participants_room_id ON room_participants(room_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_participants_room_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_participants_user_id ON room_participants(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_participants_user_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_participants_status ON room_participants(status)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_participants_status: %w", err)
	}

	// Index untuk room_messages table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_messages_room_id ON room_messages(room_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_messages_room_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_messages_sender_id ON room_messages(sender_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_messages_sender_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_room_messages_created_at ON room_messages(created_at)").Error; err != nil {
		return fmt.Errorf("failed to create idx_room_messages_created_at: %w", err)
	}

	// Index untuk user_sessions table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_user_sessions_user_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(token)").Error; err != nil {
		return fmt.Errorf("failed to create idx_user_sessions_token: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at)").Error; err != nil {
		return fmt.Errorf("failed to create idx_user_sessions_expires_at: %w", err)
	}

	// Index untuk notifications table
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_notifications_user_id: %w", err)
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read)").Error; err != nil {
		return fmt.Errorf("failed to create idx_notifications_is_read: %w", err)
	}

	logrus.Info("Database indexes created successfully")
	return nil
}

// Close menutup koneksi database
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	logrus.Info("Database connection closed")
	return nil
}

// Ping memeriksa koneksi database
func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetStats mengembalikan statistik koneksi database
func (d *Database) GetStats() map[string]interface{} {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// HealthCheck melakukan health check pada database
func (d *Database) HealthCheck() error {
	if err := d.Ping(); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Test query sederhana
	var result int
	if err := d.DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("database health check query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("database health check returned unexpected result: %d", result)
	}

	return nil
}

// getLogLevel mengkonversi string level logrus ke logger.LogLevel GORM
func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
