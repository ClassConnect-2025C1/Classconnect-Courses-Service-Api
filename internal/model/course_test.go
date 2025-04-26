package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCourseCreation(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour * 30) // 30 days in the future

	course := Course{
		ID:                  1,
		Title:               "Introduction to Go",
		Description:         "Learn the basics of Go programming language",
		CreatedBy:           "instructor@example.com",
		Capacity:            30,
		StartDate:           now,
		EndDate:             future,
		EligibilityCriteria: "Basic programming knowledge",
	}

	if course.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", course.ID)
	}
	if course.Title != "Introduction to Go" {
		t.Errorf("Expected Title to be 'Introduction to Go', got %s", course.Title)
	}
	if course.Description != "Learn the basics of Go programming language" {
		t.Errorf("Expected Description to be 'Learn the basics of Go programming language', got %s", course.Description)
	}
	if course.CreatedBy != "instructor@example.com" {
		t.Errorf("Expected CreatedBy to be 'instructor@example.com', got %s", course.CreatedBy)
	}
	if course.Capacity != 30 {
		t.Errorf("Expected Capacity to be 30, got %d", course.Capacity)
	}
	if !course.StartDate.Equal(now) {
		t.Errorf("Expected StartDate to be %v, got %v", now, course.StartDate)
	}
	if !course.EndDate.Equal(future) {
		t.Errorf("Expected EndDate to be %v, got %v", future, course.EndDate)
	}
	if course.EligibilityCriteria != "Basic programming knowledge" {
		t.Errorf("Expected EligibilityCriteria to be 'Basic programming knowledge', got %s", course.EligibilityCriteria)
	}
	if course.DeletedAt != nil {
		t.Errorf("Expected DeletedAt to be nil, got %v", course.DeletedAt)
	}
}

func TestCourseJSONMarshaling(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour * 30) // 30 days in the future

	course := Course{
		ID:                  1,
		Title:               "Introduction to Go",
		Description:         "Learn the basics of Go programming language",
		CreatedBy:           "instructor@example.com",
		Capacity:            30,
		StartDate:           now,
		EndDate:             future,
		EligibilityCriteria: "Basic programming knowledge",
	}

	// Test marshaling
	jsonData, err := json.Marshal(course)
	if err != nil {
		t.Fatalf("Failed to marshal Course to JSON: %v", err)
	}

	// Test unmarshaling
	var unmarshaledCourse Course
	err = json.Unmarshal(jsonData, &unmarshaledCourse)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to Course: %v", err)
	}

	// Verify fields were preserved
	if unmarshaledCourse.ID != course.ID {
		t.Errorf("ID not preserved during marshal/unmarshal. Expected %d, got %d", course.ID, unmarshaledCourse.ID)
	}
	if unmarshaledCourse.Title != course.Title {
		t.Errorf("Title not preserved during marshal/unmarshal. Expected %s, got %s", course.Title, unmarshaledCourse.Title)
	}
	if unmarshaledCourse.Capacity != course.Capacity {
		t.Errorf("Capacity not preserved during marshal/unmarshal. Expected %d, got %d", course.Capacity, unmarshaledCourse.Capacity)
	}
}
