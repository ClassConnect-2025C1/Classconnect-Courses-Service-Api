package model

import "time"

type CourseStatistics struct {
	ID                                      uint                      `json:"id" gorm:"primaryKey"`
	CourseID                                uint                      `json:"course_id"`
	CourseName                              string                    `json:"course_name"`
	Last10AssignmentsAverageGradeTendency   string                    `json:"last_10_days_average_grade_tendency"`
	Last10AssignmentsSubmissionRateTendency string                    `json:"last_10_days_submission_rate_tendency"`
	Suggestions                             string                    `json:"suggestions"`
	GlobalAverageGrade                      float64                   `json:"global_average_grade" gorm:"default:0"`
	GlobalSubmissionRate                    float64                   `json:"global_submission_rate" gorm:"default:0"`
	StatisticsForAssignments                []StatisticsForAssignment `json:"statistics_for_assignments" gorm:"serializer:json"`
}

type UserCourseStatistics struct {
	ID                                      uint                      `json:"id" gorm:"primaryKey"`
	CourseID                                uint                      `json:"course_id"`
	UserID                                  string                    `json:"user_id"`
	AverageGrade                            float64                   `json:"average_grade" gorm:"default:0"`
	SubmissionRate                          float64                   `json:"submission_rate" gorm:"default:0"`
	Last10AssignmentsAverageGradeTendency   string                    `json:"last_10_days_average_grade_tendency"`
	Last10AssignmentsSubmissionRateTendency string                    `json:"last_10_days_submission_rate_tendency"`
	StatisticsForAssignments                []StatisticsForAssignment `json:"statistics_for_assignments" gorm:"serializer:json"`
}

type StatisticsForAssignment struct {
	Date           time.Time `json:"date"`
	AverageGrade   float64   `json:"average_grade" gorm:"default:0"`
	SubmissionRate float64   `json:"submission_rate" gorm:"default:0"`
}
