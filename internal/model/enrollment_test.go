package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEnrollmentCreation(t *testing.T) {
	now := time.Now()
	enrollment := Enrollment{
		ID:        1,
		UserID:    "100",
		CourseID:  200,
		Role:      "instructor",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if enrollment.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", enrollment.ID)
	}
	if enrollment.UserID != "100" {
		t.Errorf("Expected UserID to be 100, got %s", enrollment.UserID)
	}
	if enrollment.CourseID != 200 {
		t.Errorf("Expected CourseID to be 200, got %d", enrollment.CourseID)
	}
	if enrollment.Role != "instructor" {
		t.Errorf("Expected Role to be 'instructor', got %s", enrollment.Role)
	}
	if enrollment.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got %s", enrollment.Email)
	}
	if enrollment.Name != "Test User" {
		t.Errorf("Expected Name to be 'Test User', got %s", enrollment.Name)
	}
}

func TestEnrollmentJSON(t *testing.T) {
	now := time.Now().Round(time.Second) // Round to avoid precision issues in comparison
	enrollment := Enrollment{
		ID:        1,
		UserID:    "100",
		CourseID:  200,
		Role:      "instructor",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test marshaling
	jsonData, err := json.Marshal(enrollment)
	if err != nil {
		t.Fatalf("Failed to marshal enrollment to JSON: %v", err)
	}

	// Test unmarshaling
	var unmarshaledEnrollment Enrollment
	err = json.Unmarshal(jsonData, &unmarshaledEnrollment)
	if err != nil {
		t.Fatalf("Failed to unmarshal enrollment from JSON: %v", err)
	}

	// Verify fields
	if unmarshaledEnrollment.ID != enrollment.ID {
		t.Errorf("ID mismatch: expected %d, got %d", enrollment.ID, unmarshaledEnrollment.ID)
	}
	if unmarshaledEnrollment.UserID != enrollment.UserID {
		t.Errorf("UserID mismatch: expected %s, got %s", enrollment.UserID, unmarshaledEnrollment.UserID)
	}
	if unmarshaledEnrollment.CourseID != enrollment.CourseID {
		t.Errorf("CourseID mismatch: expected %d, got %d", enrollment.CourseID, unmarshaledEnrollment.CourseID)
	}
	if unmarshaledEnrollment.Role != enrollment.Role {
		t.Errorf("Role mismatch: expected %s, got %s", enrollment.Role, unmarshaledEnrollment.Role)
	}
	if unmarshaledEnrollment.Email != enrollment.Email {
		t.Errorf("Email mismatch: expected %s, got %s", enrollment.Email, unmarshaledEnrollment.Email)
	}
	if unmarshaledEnrollment.Name != enrollment.Name {
		t.Errorf("Name mismatch: expected %s, got %s", enrollment.Name, unmarshaledEnrollment.Name)
	}
}
