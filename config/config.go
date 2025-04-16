package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file for local development.
// When deployed on platforms like Render, system environment variables are used.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or couldn't load it, using system environment variables (expected in deployment environments)")
	} else {
		log.Println("Loaded environment variables from .env file")
	}
}

// GetEnv obtains the value of an environment variable.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
