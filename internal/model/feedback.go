package model

import (
	"time"

	"gorm.io/gorm"
)

type CourseFeedback struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CourseID  uint           `gorm:"not null" json:"course_id"`
	UserID    string         `gorm:"not null" json:"user_id"`
	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `json:"comment"`
	Summary   string         `json:"summary"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Course Course `gorm:"foreignKey:CourseID" json:"-"`
}

type CreateFeedbackRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
	Summary string `json:"summary"`
}
