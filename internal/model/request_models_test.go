package model

import (
	"testing"
	"time"
)

func TestCreateCourseRequestToModel(t *testing.T) {
	// Create test data
	request := &CreateCourseRequest{
		Title:               "Test Course",
		Description:         "Test Description",
		CreatedBy:           "test-user",
		Capacity:            30,
		EligibilityCriteria: "None",
	}

	// Get the current time to compare with later
	before := time.Now()

	// Call the method being tested
	course := request.ToModel()

	// Get the time after the method call
	after := time.Now()

	// Verify all fields are correctly copied
	if course.Title != request.Title {
		t.Errorf("Expected Title '%s', got '%s'", request.Title, course.Title)
	}

	if course.Description != request.Description {
		t.Errorf("Expected Description '%s', got '%s'", request.Description, course.Description)
	}

	if course.CreatedBy != request.CreatedBy {
		t.Errorf("Expected CreatedBy '%s', got '%s'", request.CreatedBy, course.CreatedBy)
	}

	if course.Capacity != request.Capacity {
		t.Errorf("Expected Capacity %d, got %d", request.Capacity, course.Capacity)
	}

	if course.EligibilityCriteria != request.EligibilityCriteria {
		t.Errorf("Expected EligibilityCriteria '%s', got '%s'", request.EligibilityCriteria, course.EligibilityCriteria)
	}

	// Verify StartDate is within the expected range
	if course.StartDate.Before(before) || course.StartDate.After(after) {
		t.Errorf("StartDate %v is not between %v and %v", course.StartDate, before, after)
	}

	// Verify EndDate is approximately 4 months after StartDate
	expectedEnd := course.StartDate.AddDate(0, 4, 0)
	// Allow for small differences in time representation
	timeDiff := course.EndDate.Sub(expectedEnd)
	if timeDiff < -time.Second || timeDiff > time.Second {
		t.Errorf("Expected EndDate approximately %v, got %v (difference: %v)", expectedEnd, course.EndDate, timeDiff)
	}
}

func TestUpdateCourseRequestApplyTo(t *testing.T) {
	// Create test course with initial values
	course := &Course{
		Title:               "Original Title",
		Description:         "Original Description",
		CreatedBy:           "original-user",
		Capacity:            20,
		StartDate:           time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:             time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
		EligibilityCriteria: "Original Criteria",
	}

	// Test case 1: Update only some fields
	t.Run("PartialUpdate", func(t *testing.T) {
		// Create a copy of the original course for comparison
		origCourse := *course

		// Create update request with some fields set
		newTitle := "Updated Title"
		newCapacity := 30
		updateReq := &UpdateCourseRequest{
			Title:    &newTitle,
			Capacity: &newCapacity,
		}

		// Apply the update
		updateReq.ApplyTo(course)

		// Verify only specified fields were updated
		if course.Title != newTitle {
			t.Errorf("Expected Title '%s', got '%s'", newTitle, course.Title)
		}

		if course.Capacity != newCapacity {
			t.Errorf("Expected Capacity %d, got %d", newCapacity, course.Capacity)
		}

		// Verify unspecified fields remain unchanged
		if course.Description != origCourse.Description {
			t.Errorf("Description should not change")
		}

		if !course.StartDate.Equal(origCourse.StartDate) {
			t.Errorf("StartDate should not change")
		}

		if !course.EndDate.Equal(origCourse.EndDate) {
			t.Errorf("EndDate should not change")
		}

		if course.EligibilityCriteria != origCourse.EligibilityCriteria {
			t.Errorf("EligibilityCriteria should not change")
		}
	})

	// Test case 2: Update all fields
	t.Run("CompleteUpdate", func(t *testing.T) {
		newTitle := "New Title"
		newDesc := "New Description"
		newCapacity := 50
		newStartDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		newEndDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		newCriteria := "New Criteria"

		updateReq := &UpdateCourseRequest{
			Title:               &newTitle,
			Description:         &newDesc,
			Capacity:            &newCapacity,
			StartDate:           &newStartDate,
			EndDate:             &newEndDate,
			EligibilityCriteria: &newCriteria,
		}

		updateReq.ApplyTo(course)

		if course.Title != newTitle {
			t.Errorf("Title not updated correctly")
		}

		if course.Description != newDesc {
			t.Errorf("Description not updated correctly")
		}

		if course.Capacity != newCapacity {
			t.Errorf("Capacity not updated correctly")
		}

		if !course.StartDate.Equal(newStartDate) {
			t.Errorf("StartDate not updated correctly")
		}

		if !course.EndDate.Equal(newEndDate) {
			t.Errorf("EndDate not updated correctly")
		}

		if course.EligibilityCriteria != newCriteria {
			t.Errorf("EligibilityCriteria not updated correctly")
		}
	})

	// Test case 3: Empty update (all fields nil)
	t.Run("EmptyUpdate", func(t *testing.T) {
		origCourse := *course
		updateReq := &UpdateCourseRequest{}

		updateReq.ApplyTo(course)

		if course.Title != origCourse.Title {
			t.Errorf("Title should not change on empty update")
		}

		if course.Description != origCourse.Description {
			t.Errorf("Description should not change on empty update")
		}
	})
}
