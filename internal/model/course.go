package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Course represents a course in the system
type Course struct {
	gorm.Model
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	CreatedBy           string         `json:"created_by"`
	Capacity            int            `json:"capacity"`
	StartDate           time.Time      `json:"start_date"`
	EndDate             time.Time      `json:"end_date"`
	EligibilityCriteria pq.StringArray `json:"eligibility_criteria" gorm:"type:text[]"`
}
