package model

import (
	"time"

	"gorm.io/gorm"
)

// type File struct {
// 	ID          uint           `gorm:"primarykey" json:"id"`
// 	Name        string         `gorm:"not null" json:"name"`
// 	URL         string         `gorm:"not null" json:"url"`
// 	ContentType string         `json:"content_type"`
// 	Size        int64          `json:"size"`
// }

// Provisory: a file struct has content as binary data
type File struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Name    string `gorm:"not null" json:"name"`
	Content []byte `gorm:"not null" json:"content"`
	Size    int64  `json:"size"`
}

type Assignment struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CourseID    uint           `gorm:"not null" json:"course_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Deadline    time.Time      `json:"deadline"`
	TimeLimit   int            `json:"time_limit"` // in minutes
	Files       []File         `gorm:"many2many:assignment_files" json:"files"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Associations
	Course Course `gorm:"foreignKey:CourseID" json:"-"`
}
