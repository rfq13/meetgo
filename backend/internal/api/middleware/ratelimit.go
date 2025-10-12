package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webrtc-meeting/backend/pkg/logger"
)

// RateLimiter interface untuk rate limiting
type RateLimiter interface {
	Allow(key string) bool
	GetRemaining(key string) int
	GetResetTime(key string) time.Time
	Reset(key string)
}

// MemoryRateLimiter implementasi rate limiting di memory
type MemoryRateLimiter struct {
	limiters map[string]*TokenBucket
	mutex    sync.RWMutex
	logger   *logger.Logger
}

// TokenBucket implementasi token bucket algorithm
type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int
	lastRefill time.Time
	mutex      sync.Mutex
}

// NewMemoryRateLimiter membuat memory rate limiter baru
func NewMemoryRateLimiter(log *logger.Logger) *MemoryRateLimiter {
	limiter := &MemoryRateLimiter{
		limiters: make(map[string]*TokenBucket),
		logger:   log,
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Allow memeriksa apakah request diizinkan
func (m *MemoryRateLimiter) Allow(key string) bool {
	bucket := m.getBucket(key, 100, 10) // Default: 100 tokens, refill 10 per second
	return bucket.consume()
}

// GetRemaining mengembalikan jumlah token tersisa
func (m *MemoryRateLimiter) GetRemaining(key string) int {
	bucket := m.getBucket(key, 100, 10)
	return bucket.getTokens()
}

// GetResetTime mengembalikan waktu reset token bucket
func (m *MemoryRateLimiter) GetResetTime(key string) time.Time {
	bucket := m.getBucket(key, 100, 10)
	return bucket.lastRefill.Add(time.Second)
}

// Reset me-reset token bucket
func (m *MemoryRateLimiter) Reset(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.limiters, key)
}

// getBucket mendapatkan atau membuat token bucket baru
func (m *MemoryRateLimiter) getBucket(key string, capacity, refillRate int) *TokenBucket {
	m.mutex.RLock()
	bucket, exists := m.limiters[key]
	m.mutex.RUnlock()

	if !exists {
		m.mutex.Lock()
		// Double-check setelah dapat write lock
		bucket, exists = m.limiters[key]
		if !exists {
			bucket = NewTokenBucket(capacity, refillRate)
			m.limiters[key] = bucket
		}
		m.mutex.Unlock()
	}

	return bucket
}

// cleanup membersihkan token bucket yang tidak digunakan
func (m *MemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		now := time.Now()
		for key, bucket := range m.limiters {
			// Hapus bucket yang tidak digunakan selama 1 jam
			if now.Sub(bucket.lastRefill) > time.Hour {
				delete(m.limiters, key)
			}
		}
		m.mutex.Unlock()
	}
}

// NewTokenBucket membuat token bucket baru
func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// consume mengkonsumsi satu token
func (tb *TokenBucket) consume() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// getTokens mengembalikan jumlah token tersisa
func (tb *TokenBucket) getTokens() int {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()
	return tb.tokens
}

// refill mengisi kembali token bucket
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	tb.lastRefill = now

	// Hitung jumlah token yang harus ditambah
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	tb.tokens += tokensToAdd

	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
}

// RateLimitConfig konfigurasi rate limiting
type RateLimitConfig struct {
	// Rate limiting per menit
	RequestsPerMinute int

	// Rate limiting per jam
	RequestsPerHour int

	// Rate limiting per hari
	RequestsPerDay int

	// Custom rate limiter (jika tidak menggunakan default)
	CustomLimiter RateLimiter

	// Key generator function
	KeyGenerator func(*gin.Context) string

	// Whether to skip successful requests from counting
	SkipSuccessfulRequests bool

	// Whether to skip certain status codes from counting
	SkipStatusCodes []int
}

// DefaultRateLimitConfig konfigurasi default
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerMinute: 60,
		RequestsPerHour:   1000,
		RequestsPerDay:    10000,
		KeyGenerator:      defaultKeyGenerator,
		SkipStatusCodes:   []int{http.StatusNotFound, http.StatusBadRequest},
	}
}

// RateLimit middleware untuk rate limiting
func RateLimit(config RateLimitConfig, log *logger.Logger) gin.HandlerFunc {
	var limiter RateLimiter

	if config.CustomLimiter != nil {
		limiter = config.CustomLimiter
	} else {
		limiter = NewMemoryRateLimiter(log)
	}

	return func(c *gin.Context) {
		key := config.KeyGenerator(c)

		// Check rate limit
		if !limiter.Allow(key) {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		remaining := limiter.GetRemaining(key)
		resetTime := limiter.GetResetTime(key)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))

		// Continue processing
		c.Next()

		// Optionally skip counting successful requests
		if config.SkipSuccessfulRequests && c.Writer.Status() < http.StatusBadRequest {
			return
		}

		// Optionally skip certain status codes
		for _, code := range config.SkipStatusCodes {
			if c.Writer.Status() == code {
				return
			}
		}
	}
}

// defaultKeyGenerator generate key berdasarkan IP address
func defaultKeyGenerator(c *gin.Context) string {
	return c.ClientIP()
}

// UserBasedKeyGenerator generate key berdasarkan user ID (untuk authenticated requests)
func UserBasedKeyGenerator(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(string); ok {
			return "user:" + uid
		}
	}
	return "ip:" + c.ClientIP()
}

// EndpointBasedKeyGenerator generate key berdasarkan endpoint
func EndpointBasedKeyGenerator(c *gin.Context) string {
	return fmt.Sprintf("%s:%s:%s", c.ClientIP(), c.Request.Method, c.FullPath())
}

// SimpleRateLimit middleware rate limiting sederhana
func SimpleRateLimit(requestsPerMinute int, log *logger.Logger) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerMinute: requestsPerMinute,
		KeyGenerator:      defaultKeyGenerator,
	}
	return RateLimit(config, log)
}

// UserRateLimit middleware rate limiting berdasarkan user
func UserRateLimit(requestsPerMinute int, log *logger.Logger) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerMinute: requestsPerMinute,
		KeyGenerator:      UserBasedKeyGenerator,
	}
	return RateLimit(config, log)
}

// APIRateLimit middleware rate limiting untuk API endpoints
func APIRateLimit(log *logger.Logger) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerMinute: 100,
		RequestsPerHour:   2000,
		RequestsPerDay:    50000,
		KeyGenerator:      EndpointBasedKeyGenerator,
		SkipStatusCodes:   []int{http.StatusNotFound, http.StatusBadRequest, http.StatusUnauthorized},
	}
	return RateLimit(config, log)
}

// AuthRateLimit middleware rate limiting untuk auth endpoints
func AuthRateLimit(log *logger.Logger) gin.HandlerFunc {
	config := RateLimitConfig{
		RequestsPerMinute: 10, // Lebih strict untuk auth
		RequestsPerHour:   100,
		RequestsPerDay:    500,
		KeyGenerator:      defaultKeyGenerator,
		SkipStatusCodes:   []int{http.StatusNotFound},
	}
	return RateLimit(config, log)
}
