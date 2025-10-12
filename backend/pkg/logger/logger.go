package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger struct untuk menyimpan logger instance
type Logger struct {
	*logrus.Logger
}

// Config konfigurasi logger
type Config struct {
	Level      string
	Format     string // json atau text
	Output     string // stdout, stderr, atau file path
	MaxSize    int    // max size dalam MB untuk rotation
	MaxBackups int    // max backup files
	MaxAge     int    // max age dalam hari
	Compress   bool   // compress backup files
}

// NewLogger membuat logger instance baru
func NewLogger(cfg *Config) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
		logger.Warnf("Invalid log level '%s', using 'info'", cfg.Level)
	}
	logger.SetLevel(level)

	// Set formatter
	switch cfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
				logrus.FieldKeyFile:  "file",
			},
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
			ForceColors:     true,
		})
	}

	// Set output
	var output io.Writer
	switch cfg.Output {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		// Jika output adalah file path, buat file writer
		if cfg.Output != "" {
			// Buat direktori jika belum ada
			dir := filepath.Dir(cfg.Output)
			if err := os.MkdirAll(dir, 0755); err != nil {
				logger.Errorf("Failed to create log directory: %v", err)
				output = os.Stdout
			} else {
				file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logger.Errorf("Failed to open log file: %v", err)
					output = os.Stdout
				} else {
					output = file
				}
			}
		} else {
			output = os.Stdout
		}
	}
	logger.SetOutput(output)

	return &Logger{Logger: logger}
}

// WithRequestID menambahkan request ID ke log entry
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField("request_id", requestID)
}

// WithUserID menambahkan user ID ke log entry
func (l *Logger) WithUserID(userID string) *logrus.Entry {
	return l.WithField("user_id", userID)
}

// WithRoomID menambahkan room ID ke log entry
func (l *Logger) WithRoomID(roomID string) *logrus.Entry {
	return l.WithField("room_id", roomID)
}

// WithError menambahkan error ke log entry
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithFields menambahkan multiple fields ke log entry
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// LogRequest log HTTP request
func (l *Logger) LogRequest(method, path string, statusCode int, duration time.Duration, clientIP string) {
	l.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration.Milliseconds(),
		"client_ip":   clientIP,
	}).Info("HTTP Request")
}

// LogError log error dengan context
func (l *Logger) LogError(err error, message string, fields ...logrus.Fields) {
	entry := l.WithError(err)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(message)
}

// LogPanic log panic dan recovery
func (l *Logger) LogPanic(recovered interface{}, message string) {
	l.WithField("panic", recovered).Error(message)
}

// LogDatabaseQuery log database query
func (l *Logger) LogDatabaseQuery(query string, duration time.Duration, rowsAffected int64) {
	l.WithFields(logrus.Fields{
		"query":         query,
		"duration":      duration.Milliseconds(),
		"rows_affected": rowsAffected,
	}).Debug("Database Query")
}

// LogWebRTCEvent log WebRTC events
func (l *Logger) LogWebRTCEvent(eventType, roomID, userID string, data interface{}) {
	l.WithFields(logrus.Fields{
		"event_type": eventType,
		"room_id":    roomID,
		"user_id":    userID,
		"data":       data,
	}).Info("WebRTC Event")
}

// LogAuthEvent log authentication events
func (l *Logger) LogAuthEvent(action, userID, clientIP string, success bool) {
	l.WithFields(logrus.Fields{
		"action":    action,
		"user_id":   userID,
		"client_ip": clientIP,
		"success":   success,
	}).Info("Authentication Event")
}

// LogBusinessEvent log business logic events
func (l *Logger) LogBusinessEvent(event string, fields logrus.Fields) {
	l.WithFields(fields).Infof("Business Event: %s", event)
}

// GetLogger mengembalikan logger instance global
var defaultLogger *Logger

// InitLogger menginisialisasi logger global
func InitLogger(cfg *Config) {
	defaultLogger = NewLogger(cfg)
}

// GetDefaultLogger mengembalikan default logger
func GetDefaultLogger() *Logger {
	if defaultLogger == nil {
		// Jika belum diinisialisasi, buat default logger
		defaultLogger = NewLogger(&Config{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		})
	}
	return defaultLogger
}

// Helper functions untuk global logger
func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return GetDefaultLogger().WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetDefaultLogger().WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return GetDefaultLogger().WithError(err)
}

// GinMiddleware mengembalikan Gin middleware untuk logging
func GinMiddleware(logger *Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.LogRequest(
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
		)
		return ""
	})
}

// RecoveryMiddleware mengembalikan Gin middleware untuk recovery
func RecoveryMiddleware(logger *Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.LogPanic(recovered, "Panic recovered")
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}
