package models

import "time"

// Enrollment representa la relación entre un usuario y un curso en la db
type Enrollment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	CourseID  uint      `json:"course_id" gorm:"index"`
	Role      string    `json:"role" gorm:"default:student"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
