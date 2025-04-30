package model

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CourseID    uint           `gorm:"not null" json:"course_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Deadline    time.Time      `json:"deadline"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Course Course `gorm:"foreignKey:CourseID" json:"-"`
}
