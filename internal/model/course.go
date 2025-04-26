package model

import "time"

// Course representa la estructura de un curso
type Course struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	Title               string     `json:"title" binding:"required"`
	Description         string     `json:"description"`
	CreatedBy           string     `json:"created_by" binding:"required"`
	Capacity            int        `json:"capacity" binding:"required,gte=1"`
	StartDate           time.Time  `json:"start_date" binding:"required"`
	EndDate             time.Time  `json:"end_date" binding:"required"`
	EligibilityCriteria string     `json:"eligibility_criteria"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}
