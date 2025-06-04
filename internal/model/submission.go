package model

import (
	"time"
)

type Submission struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	CourseID     uint             `json:"course_id" gorm:"not null"`
	AssignmentID uint             `json:"assignment_id" gorm:"not null"`
	UserID       string           `json:"user_id" gorm:"not null"`
	Content      string           `json:"content" gorm:"not null"`
	SubmittedAt  time.Time        `json:"submitted_at" gorm:"autoCreateTime"`
	Grade        uint             `json:"grade" gorm:"check:grade >= 0 AND grade <= 100"`
	Feedback     string           `json:"feedback"`
	Files        []SubmissionFile `gorm:"many2many:submission_files_join" json:"files"`
}

type SubmissionFile struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Content []byte `gorm:"not null" json:"content"`
	Size    int64  `json:"size"`
}
