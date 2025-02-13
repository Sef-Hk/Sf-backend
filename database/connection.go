package database

import (
	"fmt"
	"log"
	"os"

	"github.com/Sef-Hk/State-Full-authentication/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ Could not load .env file, using system environment variables")
	}
	host := os.Getenv("PGHOST")
	port := os.Getenv("PGPORT")
	user := os.Getenv("PGUSER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("PGDATABASE")
	// Use PGSSLMODE environment variable with a default value of "require"
	sslmode := os.Getenv("PGSSLMODE")
	if sslmode == "" {
		sslmode = "require"
	}
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Database connection variables are missing!")
	}
	// Construct DSN with the sslmode variable
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Could not connect to the database: %v", err)
	}

	DB = connection
	connection.AutoMigrate(&models.User{})
	log.Println("✅ Successfully connected to the database!")
}
