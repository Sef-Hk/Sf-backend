package config

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

// Global session store
var Store *session.Store

func InitSessionStore() {
	if Store != nil {
		return // Prevent multiple initializations
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("❌ Missing REDIS_URL in .env file")
	}

	redisStore := redis.New(redis.Config{
		URL:       redisURL,
		TLSConfig: &tls.Config{},
	})

	Store = session.New(session.Config{
		Storage:        redisStore,
		Expiration:     24 * time.Hour,
		CookieSecure:   true, // Set to true in production
		CookieHTTPOnly: true,
		CookieSameSite: "None",
		CookiePath:     "/",
	})
	log.Println("✅ Session store initialized with Redis!")
}
