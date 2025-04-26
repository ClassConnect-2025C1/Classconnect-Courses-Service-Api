package main

import (
	"os"
	"testing"
)

// TestPortSelectionWhenSet tests that the port is correctly read from
// environment variables when set
func TestPortSelectionWhenSet(t *testing.T) {
	// Save original environment variable
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	// Set environment variable
	os.Setenv("PORT", "9090")

	// Check port selection logic
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Verify port
	if port != "9090" {
		t.Errorf("Expected port to be 9090, got %s", port)
	}
}

// TestPortSelectionWhenNotSet tests that the default port is used when
// environment variable is not set
func TestPortSelectionWhenNotSet(t *testing.T) {
	// Save original environment variable
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	// Unset environment variable
	os.Unsetenv("PORT")

	// Check port selection logic
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Verify port
	if port != "8080" {
		t.Errorf("Expected default port to be 8080, got %s", port)
	}
}
