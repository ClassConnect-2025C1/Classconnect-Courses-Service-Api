package model

import "time"

type CourseStatistics struct {
	CourseID             uint                `json:"course_id"`
	CourseName           string              `json:"course_name"`
	GlobalAverageGrade   float64             `json:"global_average_grade"`
	GlobalSubmissionRate float64             `json:"global_submission_rate"`
	StatisticsForDates   []StatisticsForDate `json:"statistics_for_dates"`
}

type UserCourseStatistics struct {
	AverageGrade       float64             `json:"average_grade"`
	SubmissionRate     float64             `json:"submission_rate"`
	StatisticsForDates []StatisticsForDate `json:"statistics_for_dates"`
}

type StatisticsForDate struct {
	Date           time.Time `json:"date"`
	AverageGrade   float64   `json:"average_grade"`
	SubmissionRate float64   `json:"submission_rate"`
}
