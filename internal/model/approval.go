package model

import (
	"gorm.io/gorm"
)

// CourseApproval represents a course approval for a user
type CourseApproval struct {
	gorm.Model
	UserID     string `gorm:"index" json:"user_id"`
	CourseID   uint   `gorm:"index" json:"course_id"`
	CourseName string `json:"course_name"`
}
