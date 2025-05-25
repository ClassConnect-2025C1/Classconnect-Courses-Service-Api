package model

import "time"

// CreateCourseRequest represents the input for creating a course
type CreateCourseRequest struct {
	Title               string   `json:"title" binding:"required"`
	Description         string   `json:"description"`
	CreatedBy           string   `json:"created_by" binding:"required"`
	Capacity            int      `json:"capacity" binding:"required,gte=1"`
	EligibilityCriteria []string `json:"eligibility_criteria"`
	TeachingAssistants  []string `json:"teaching_assistants"`
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
		TeachingAssistants:  r.TeachingAssistants,
	}
}

// UpdateCourseRequest represents the input for updating a course
type UpdateCourseRequest struct {
	Title               *string    `json:"title"`
	Description         *string    `json:"description"`
	Capacity            *int       `json:"capacity"`
	StartDate           *time.Time `json:"start_date"`
	EndDate             *time.Time `json:"end_date"`
	EligibilityCriteria *[]string  `json:"eligibility_criteria"`
	TeachingAssistants  *[]string  `json:"teaching_assistants"`
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
	if r.TeachingAssistants != nil {
		course.TeachingAssistants = *r.TeachingAssistants // Apply TA updates
	}
}

type CreateAssignmentRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline" binding:"required"`
	TimeLimit   int       `json:"time_limit"` // in minutes
	Files       []File    `json:"files"`      // Provisory: a file struct has content as binary data
}

type UpdateAssignmentRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	TimeLimit   int       `json:"time_limit"` // in minutes
	Files       []File    `json:"files"`      // Provisory: a file struct has content as binary data
}

type CreateSubmissionRequest struct {
	CourseID     uint             `json:"course_id" binding:"required"`
	AssignmentID uint             `json:"assignment_id" binding:"required"`
	Content      string           `json:"content" binding:"required"`
	Files        []SubmissionFile `json:"files"`
}

type GradeSubmissionRequest struct {
	Grade    uint   `json:"grade" binding:"required,gte=0,lte=100"`
	Feedback string `json:"feedback"`
}
