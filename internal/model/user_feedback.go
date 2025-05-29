package model

import (
	"gorm.io/gorm"
)

// UserFeedback represents feedback given to a student by an instructor
type UserFeedback struct {
	gorm.Model
	CourseID     uint   `json:"course_id"`
	Course       Course `json:"course" gorm:"foreignKey:CourseID"`
	StudentID    string `json:"student_id"`    // The user receiving feedback
	InstructorID string `json:"instructor_id"` // The user giving feedback
	Comment      string `json:"comment"`
	Rating       uint   `json:"rating"` // A numeric rating, e.g., 1-5
}
