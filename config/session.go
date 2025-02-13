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
		CookieSameSite: "Lax",
		CookiePath:     "/",
	})
	log.Println("✅ Session store initialized with Redis!")
}

// package config

// import (
// 	"crypto/tls"
// 	"log"
// 	"time"

// 	"github.com/gofiber/fiber/v2/middleware/session"
// 	"github.com/gofiber/storage/redis"
// )

// func InitSessionStore() *session.Store {
// 	// Initialize Redis storage
// 	redisStore := redis.New(redis.Config{
// 		Host:      "allowed-goldfish-14767.upstash.io",                          // Redis host
// 		Port:      6379,                                                         // Redis port
// 		Password:  "ATmvAAIjcDExYjVhNzQyZmQ2ZmI0OTVjOTBhMmRhN2IxYmJmNTkxNnAxMA", // Redis password
// 		Database:  0,
// 		TLSConfig: &tls.Config{},
// 	})

// 	// Initialize session store with Redis as storage
// 	store := session.New(session.Config{
// 		Storage:        redisStore, // ✅ Use Redis for session storage
// 		Expiration:     24 * time.Hour,
// 		CookieSecure:   true, // Set to true in production
// 		CookieHTTPOnly: true,
// 		CookieSameSite: "Lax",
// 		CookiePath:     "/",
// 	})

// 	log.Println("✅ Session store initialized with Redis!")
// 	return store
// }

// // Global variable to hold session store
// var Store *session.Store

// func init() {
// 	Store = InitSessionStore()
// }
