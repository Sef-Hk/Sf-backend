package middleware

import (
	"time"

	"github.com/Sef-Hk/State-Full-authentication/backend/config"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err != nil {
		_ = sess.Destroy() // Ensure session is cleared if retrieval fails
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - No session found",
		})
	}

	userID := sess.Get("user_id")
	if userID == nil {
		_ = sess.Destroy() // Clear session if user_id is missing
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Invalid session",
		})
	}

	// Retrieve last activity timestamp
	lastActivity := sess.Get("last_activity")
	if lastActivity != nil {
		lastActivityTimestamp, ok := lastActivity.(int64) // Use int64 for Redis compatibility
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Session timestamp error",
			})
		}

		// Convert timestamp back to time.Time
		lastActivityTime := time.Unix(lastActivityTimestamp, 0)

		// Destroy the session if inactive for more than 30 minutes
		if time.Since(lastActivityTime) > 30*time.Minute {
			_ = sess.Destroy()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session expired due to inactivity",
			})
		}
	}

	// Update last activity timestamp
	sess.Set("last_activity", time.Now().Unix()) // Store as Unix timestamp for Redis
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update session timestamp",
		})
	}

	return c.Next()
}
