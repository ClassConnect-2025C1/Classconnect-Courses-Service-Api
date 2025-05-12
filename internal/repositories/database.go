package repositories

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"templateGo/internal/model"
	"time"

	"github.com/lib/pq"
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
	password := getEnv("DB_PASSWORD", "1234")
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

	// Register array types
	DB.Callback().Create().Before("gorm:create").Register("pq_array_handler", arrayHandlerCreate)
	DB.Callback().Update().Before("gorm:update").Register("pq_array_handler", arrayHandlerUpdate)

	// Auto migrate model
	if err := DB.AutoMigrate(&model.Course{}, &model.Enrollment{}, &model.CourseFeedback{}, &model.Assignment{}, &model.CourseApproval{}, &model.Submission{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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

func arrayHandlerCreate(db *gorm.DB) {
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.Fields {
			if field.FieldType.Kind() == reflect.Slice && field.FieldType.Elem().Kind() == reflect.String {
				if v, ok := db.Statement.ReflectValue.FieldByName(field.Name).Interface().([]string); ok {
					db.Statement.SetColumn(field.DBName, pq.Array(v))
				}
			}
		}
	}
}

func arrayHandlerUpdate(db *gorm.DB) {
	// Similar to arrayHandlerCreate but for updates
	// ...
	if db.Statement.Schema != nil {
		for _, field := range db.Statement.Schema.Fields {
			if field.FieldType.Kind() == reflect.Slice && field.FieldType.Elem().Kind() == reflect.String {
				if v, ok := db.Statement.ReflectValue.FieldByName(field.Name).Interface().([]string); ok {
					db.Statement.SetColumn(field.DBName, pq.Array(v))
				}
			}
		}
	}
}
