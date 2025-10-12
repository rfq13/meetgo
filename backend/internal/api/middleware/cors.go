package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig konfigurasi CORS
type CORSConfig struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposedHeaders     []string
	AllowCredentials   bool
	MaxAge             int
	OptionsPassthrough bool
}

// DefaultCORSConfig mengembalikan konfigurasi CORS default
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Accept",
			"Cache-Control",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// ProductionCORSConfig mengembalikan konfigurasi CORS untuk production
func ProductionCORSConfig(allowedDomains []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins:   allowedDomains,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "Accept"},
		ExposedHeaders:   []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}

// CORS middleware untuk Cross-Origin Resource Sharing
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", getAllowedOrigin(config, origin))
		c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			if config.OptionsPassthrough {
				c.Next()
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// getAllowedOrigin memeriksa apakah origin diizinkan
func getAllowedOrigin(config CORSConfig, origin string) string {
	if origin == "" {
		return ""
	}

	// Jika wildcard (*) digunakan
	for _, allowedOrigin := range config.AllowedOrigins {
		if allowedOrigin == "*" {
			return "*"
		}
		if allowedOrigin == origin {
			return origin
		}
	}

	// Jika origin tidak diizinkan, return empty string atau first allowed origin
	if len(config.AllowedOrigins) > 0 {
		return config.AllowedOrigins[0]
	}

	return ""
}

// SimpleCORS middleware CORS sederhana untuk development
func SimpleCORS() gin.HandlerFunc {
	return CORS(DefaultCORSConfig())
}

// RestrictiveCORS middleware CORS yang lebih restriktif
func RestrictiveCORS(allowedOrigins []string) gin.HandlerFunc {
	config := CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           3600, // 1 hour
	}
	return CORS(config)
}
