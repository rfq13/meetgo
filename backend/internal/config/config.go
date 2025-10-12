package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config struct untuk menyimpan semua konfigurasi aplikasi
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Janus    JanusConfig
	Email    EmailConfig
	Logger   LoggerConfig
}

// ServerConfig konfigurasi server
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig konfigurasi database
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig konfigurasi Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig konfigurasi JWT
type JWTConfig struct {
	Secret         string
	RefreshSecret  string
	ExpirationTime time.Duration
	RefreshTime    time.Duration
}

// JanusConfig konfigurasi Janus WebRTC server
type JanusConfig struct {
	WebSocketURL string
	HTTPURL      string
	AdminSecret  string
	APISecret    string
}

// EmailConfig konfigurasi email
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	From         string
}

// LoggerConfig konfigurasi logger
type LoggerConfig struct {
	Level  string
	Format string
}

// LoadConfig memuat konfigurasi dari environment variables
func LoadConfig() (*Config, error) {
	// Load .env file jika ada
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "webrtc_meeting"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "your-secret-key"),
			RefreshSecret:  getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			ExpirationTime: getDurationEnv("JWT_EXPIRATION_TIME", 24*time.Hour),
			RefreshTime:    getDurationEnv("JWT_REFRESH_TIME", 168*time.Hour), // 7 days
		},
		Janus: JanusConfig{
			WebSocketURL: getEnv("JANUS_WS_URL", "ws://localhost:8188"),
			HTTPURL:      getEnv("JANUS_HTTP_URL", "http://localhost:8088/janus"),
			AdminSecret:  getEnv("JANUS_ADMIN_SECRET", "janusrocks"),
			APISecret:    getEnv("JANUS_API_SECRET", "janusrocks"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getIntEnv("SMTP_PORT", 587),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			From:         getEnv("EMAIL_FROM", "noreply@webrtc-meeting.com"),
		},
		Logger: LoggerConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validasi konfigurasi yang diperlukan
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// Validate memvalidasi konfigurasi
func (c *Config) Validate() error {
	if c.JWT.Secret == "your-secret-key" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}

	if c.JWT.RefreshSecret == "your-refresh-secret-key" {
		return fmt.Errorf("JWT_REFRESH_SECRET must be set in production")
	}

	return nil
}

// GetDatabaseURL mengembalikan database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisURL mengembalikan Redis connection string
func (c *Config) GetRedisURL() string {
	if c.Redis.Password != "" {
		return fmt.Sprintf("%s:%s@%s:%s/%d",
			c.Redis.Password,
			c.Redis.Password,
			c.Redis.Host,
			c.Redis.Port,
			c.Redis.DB,
		)
	}
	return fmt.Sprintf("%s:%s/%d",
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.DB,
	)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
