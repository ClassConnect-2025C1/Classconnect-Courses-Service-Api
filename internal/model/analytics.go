package model

// CourseAnalytics represents GENERAL analytics data for a course
// Contains course-wide statistics like total enrollment, average grades across all students, etc.
// Statistics are stored as JSON string for maximum flexibility
type CourseAnalytics struct {
	CourseID   uint   `json:"course_id" gorm:"primaryKey"`
	Statistics []byte `json:"statistics" gorm:"type:json"` // JSON with course-wide statistics

	// Foreign key relationship
	Course Course `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// UserCourseAnalytics represents INDIVIDUAL analytics data for a specific user enrolled in a course
// Contains user-specific statistics like personal grades, individual progress, time spent, etc.
// This is completely different from CourseAnalytics - it tracks individual performance, not course-wide metrics
type UserCourseAnalytics struct {
	UserID     string `json:"user_id" gorm:"primaryKey"`   // User identifier
	CourseID   uint   `json:"course_id" gorm:"primaryKey"` // Course ID
	Statistics []byte `json:"statistics" gorm:"type:json"` // JSON with course-wide statistics

	// Foreign key relationship
	Course Course `json:"course" gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
