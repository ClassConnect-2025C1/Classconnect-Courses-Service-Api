package model

import "time"

type CourseStatistics struct {
	CourseID                         uint                `json:"course_id"`
	CourseName                       string              `json:"course_name"`
	Last10DaysAverageGradeTendency   string              `json:"last_10_days_average_grade_tendency"`
	Last10DaysSubmissionRateTendency string              `json:"last_10_days_submission_rate_tendency"`
	GlobalAverageGrade               float64             `json:"global_average_grade" gorm:"default:0"`
	GlobalSubmissionRate             float64             `json:"global_submission_rate" gorm:"default:0"`
	StatisticsForDates               []StatisticsForDate `json:"statistics_for_dates"`
}

type UserCourseStatistics struct {
	CourseID                         uint                `json:"course_id"`
	AverageGrade                     float64             `json:"average_grade" gorm:"default:0"`
	SubmissionRate                   float64             `json:"submission_rate" gorm:"default:0"`
	Last10DaysAverageGradeTendency   string              `json:"last_10_days_average_grade_tendency"`
	Last10DaysSubmissionRateTendency string              `json:"last_10_days_submission_rate_tendency"`
	StatisticsForDates               []StatisticsForDate `json:"statistics_for_dates"`
}

type StatisticsForDate struct {
	Date           time.Time `json:"date"`
	AverageGrade   float64   `json:"average_grade" gorm:"default:0"`
	SubmissionRate float64   `json:"submission_rate" gorm:"default:0"`
}
