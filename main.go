package main

import (
	"log"
	"os"

	"github.com/Sef-Hk/State-Full-authentication/backend/config"
	"github.com/Sef-Hk/State-Full-authentication/backend/database"
	"github.com/Sef-Hk/State-Full-authentication/backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to the database
	database.Connect()

	// Initialize session store once
	config.InitSessionStore()
	allowedOrigin := os.Getenv("ALLOW_ORIGIN")
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigin,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.SetUp(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local dev
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
