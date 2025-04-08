package sql

import (
	"fmt"
	"log"
	"os"
	"templateGo/internals/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB instance
var DB *gorm.DB

// ConnectDB establishes connection to the database
func ConnectDB() error {
	var err error

	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	if err := DB.AutoMigrate(&models.Course{}); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	// Get the underlying SQL DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Connected to SQL Database")
	return nil
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return fmt.Errorf("error getting SQL DB: %w", err)
		}

		err = sqlDB.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}

		fmt.Println("Database connection closed")
		DB = nil
	}
	return nil
}
