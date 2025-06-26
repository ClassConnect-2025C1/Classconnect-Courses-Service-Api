package model

import (
	"gorm.io/gorm"
)

// CourseAnalytics represents GENERAL analytics data for a course
// Contains course-wide statistics like total enrollment, average grades across all students, etc.
// Statistics are stored as JSON string for maximum flexibility
type CourseAnalytics struct {
	gorm.Model
	CourseID   uint   `json:"course_id" gorm:"not null;index"`
	Statistics string `json:"statistics" gorm:"type:text"` // JSON with course-wide statistics

	// Foreign key relationship
	Course Course `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// UserCourseAnalytics represents INDIVIDUAL analytics data for a specific user enrolled in a course
// Contains user-specific statistics like personal grades, individual progress, time spent, etc.
// This is completely different from CourseAnalytics - it tracks individual performance, not course-wide metrics
type UserCourseAnalytics struct {
	gorm.Model
	UserID     string `json:"user_id" gorm:"not null;index:idx_user_course,unique"`   // User identifier
	CourseID   uint   `json:"course_id" gorm:"not null;index:idx_user_course,unique"` // Course ID
	Statistics string `json:"statistics" gorm:"type:text"`                            // JSON with user-specific statistics

	// Foreign key relationship
	Course Course `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
