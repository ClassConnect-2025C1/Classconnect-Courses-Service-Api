package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test case 1: Environment variable exists
	key := "TEST_EXISTING_ENV"
	expectedValue := "test_value"
	os.Setenv(key, expectedValue)
	defer os.Unsetenv(key)

	result := GetEnv(key, "default_value")
	if result != expectedValue {
		t.Errorf("GetEnv(%s, \"default_value\") = %s; want %s", key, result, expectedValue)
	}

	// Test case 2: Environment variable doesn't exist
	nonExistentKey := "TEST_NON_EXISTENT_ENV"
	defaultValue := "default_value"
	os.Unsetenv(nonExistentKey) // Make sure it doesn't exist

	result = GetEnv(nonExistentKey, defaultValue)
	if result != defaultValue {
		t.Errorf("GetEnv(%s, %s) = %s; want %s", nonExistentKey, defaultValue, result, defaultValue)
	}
}

func TestLoadEnv(t *testing.T) {
	// This is a simple test to ensure LoadEnv doesn't crash
	// when no .env file exists
	LoadEnv()
	// If we get here without panicking, the test passes
}
