package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Room-Balance/room-balance-backend.git/db/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Create connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connection established successfully.")
}

func MigrateDB() {
	// Ensure ENUM types exist before migration
	createEnums()

	// Check if migration is needed
	if isDatabaseSetup() {
		log.Println("Database already set up. Skipping migration.")
		return
	}

	err := DB.AutoMigrate(
		&models.User{},
		&models.House{},
		&models.UserHouse{},
		&models.Task{},
		&models.Expense{},
		&models.Event{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed!")
}

// createEnums ensures that required ENUM types are created if not already present
func createEnums() {
	enumQueries := map[string]string{
		"task_type":      `CREATE TYPE task_type AS ENUM ('one-time', 'repetitive');`,
		"task_status":    `CREATE TYPE task_status AS ENUM ('pending', 'in-progress', 'done');`,
		"task_frequency": `CREATE TYPE task_frequency AS ENUM ('daily', 'weekly', 'monthly');`,
	}

	for enumType, query := range enumQueries {
		if !enumExists(enumType) {
			err := DB.Exec(query).Error
			if err != nil {
				log.Fatalf("Failed to create ENUM type %s: %v", enumType, err)
			}
			log.Printf("ENUM type %s created successfully.", enumType)
		} else {
			log.Printf("ENUM type %s already exists. Skipping creation.", enumType)
		}
	}
}

// enumExists checks if a given ENUM type already exists in the database
func enumExists(enumName string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?);`
	err := DB.Raw(query, enumName).Scan(&exists).Error
	if err != nil {
		log.Fatalf("Failed to check existence of ENUM type %s: %v", enumName, err)
	}
	return exists
}

// isDatabaseSetup checks if essential tables are already present
func isDatabaseSetup() bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');`
	err := DB.Raw(query).Scan(&exists).Error
	if err != nil {
		log.Fatalf("Failed to check database setup: %v", err)
	}
	return exists
}
