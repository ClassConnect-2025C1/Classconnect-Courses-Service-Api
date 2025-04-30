package model

import "time"

// CreateCourseRequest represents the input for creating a course
type CreateCourseRequest struct {
	Title               string `json:"title" binding:"required"`
	Description         string `json:"description"`
	CreatedBy           string `json:"created_by" binding:"required"`
	Capacity            int    `json:"capacity" binding:"required,gte=1"`
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

type CreateAssignmentRequest struct {
	CourseID    uint      `json:"course_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

func (r *CreateAssignmentRequest) ToModel() *Assignment {
	return &Assignment{
		CourseID:    r.CourseID,
		Title:       r.Title,
		Description: r.Description,
		Deadline:    r.Deadline,
		CreatedAt:   time.Now(),
	}
}

type UpdateAssignmentRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
}

func (r *UpdateAssignmentRequest) ApplyTo(assignment *Assignment) {
	assignment.Title = r.Title
	assignment.Description = r.Description
	assignment.Deadline = r.Deadline
}
