package models

import "time"

// CreateCourseRequest represents the input for creating a course
type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,gte=1"`
	// dejamos por dafult 4 meses
	// StartDate           string `json:"start_date" binding:"required"`
	// EndDate             string `json:"end_date" binding:"required"`
	EligibilityCriteria string `json:"eligibility_criteria"`
}

// ToModel converts API request to internal Course model
func (r *CreateCourseRequest) ToModel() *Course {
	return &Course{
		Title:               r.Title,
		Description:         r.Description,
		CreatedBy:           r.CreatedBy,
		Capacity:            r.Capacity,
		StartDate:           time.Now(),
		EndDate:             time.Now().AddDate(0, 4, 0), // 4 months by default
		EligibilityCriteria: r.EligibilityCriteria,
	}
}

// UpdateCourseRequest represents the input for updating a course
type UpdateCourseRequest struct {
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Capacity            *int       `json:"capacity"`
	StartDate           *time.Time `json:"start_date"`
	EndDate             *time.Time `json:"end_date"`
	EligibilityCriteria *string    `json:"eligibility_criteria"`
}

// ApplyTo applies the update request to an existing course
func (r *UpdateCourseRequest) ApplyTo(course *Course) {
	if r.Title != nil {
		course.Title = *r.Title
	}
	if r.Description != nil {
		course.Description = *r.Description
	}
	if r.Capacity != nil {
		course.Capacity = *r.Capacity
	}
	if r.StartDate != nil {
		course.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		course.EndDate = *r.EndDate
	}
	if r.EligibilityCriteria != nil {
		course.EligibilityCriteria = *r.EligibilityCriteria
	}
}
