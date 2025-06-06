package model

import "time"

type CourseStatistics struct {
	CourseID             uint                `json:"course_id"`
	CourseName           string              `json:"course_name"`
	GlobalAverageGrade   float64             `json:"global_average_grade" gorm:"default:0"`
	GlobalSubmissionRate float64             `json:"global_submission_rate" gorm:"default:0"`
	StatisticsForDates   []StatisticsForDate `json:"statistics_for_dates"`
}

type UserCourseStatistics struct {
	AverageGrade       float64             `json:"average_grade" gorm:"default:0"`
	SubmissionRate     float64             `json:"submission_rate" gorm:"default:0"`
	StatisticsForDates []StatisticsForDate `json:"statistics_for_dates"`
}

type StatisticsForDate struct {
	Date           time.Time `json:"date"`
	AverageGrade   float64   `json:"average_grade" gorm:"default:0"`
	SubmissionRate float64   `json:"submission_rate" gorm:"default:0"`
}
