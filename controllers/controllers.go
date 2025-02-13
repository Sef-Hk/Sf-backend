package controllers

import (
	// Assuming you have a database package

	"log"
	"time"

	"github.com/Sef-Hk/State-Full-authentication/backend/config"
	"github.com/Sef-Hk/State-Full-authentication/backend/database"
	"github.com/Sef-Hk/State-Full-authentication/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	// Create a struct for parsing user input separately
	type RegisterInput struct {
		FullName    string         `json:"full_name"`
		Email       string         `json:"email"`
		Password    string         `json:"password"` // Accept password as a string in JSON
		PhoneNumber string         `json:"phone_number"`
		DateOfBirth string         `json:"date_of_birth"`
		City        string         `json:"city"`
		Country     string         `json:"country"`
		Role        string         `json:"role"`
		Skills      pq.StringArray `json:"skills"`
	}

	var input RegisterInput

	// Parse request body into the temporary struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid user structure",
			"details": err.Error(),
		})
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email is already registered"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not process password",
		})
	}

	// Create the user with the hashed password
	user := models.User{
		FullName:    input.FullName,
		Email:       input.Email,
		Password:    hashedPassword, // Store as []byte
		PhoneNumber: input.PhoneNumber,
		DateOfBirth: input.DateOfBirth,
		City:        input.City,
		Country:     input.Country,
		Role:        input.Role,
		Skills:      input.Skills,
	}

	// Insert into the database
	result := database.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to register user",
			"details": result.Error.Error(),
		})
	}

	// Remove password before sending response
	user.Password = nil

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func Login(c *fiber.Ctx) error {
	// Request body struct
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Compare stored password hash with provided password
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate a new session
	sess, err := config.Store.Get(c)
	if err != nil {
		log.Println("Session error:", err) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create session",
		})
	}

	// Store session data
	sess.Set("user_id", user.ID)
	sess.Set("user_role", user.Role)
	sess.Set("last_activity", time.Now().Unix()) // ✅ Store as Unix timestamp for Redis compatibility

	if err := sess.Save(); err != nil {
		log.Println("Session save error:", err) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not save session",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user": fiber.Map{
			"id":       user.ID,
			"fullName": user.FullName,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized - No session found"})
	}

	userID := sess.Get("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Fetch user details
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":            user.ID,
		"full_name":     user.FullName,
		"email":         user.Email,
		"role":          user.Role,
		"skills":        user.Skills,
		"phone_number":  user.PhoneNumber,
		"country":       user.Country,
		"city":          user.City,
		"date_of_birth": user.DateOfBirth,
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := config.Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Session retrieval error"})
	}

	// Destroy session
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to logout"})
	}

	return c.JSON(fiber.Map{"message": "Logout successful"})
}
