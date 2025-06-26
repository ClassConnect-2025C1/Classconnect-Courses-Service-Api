package main

import (
	"encoding/json"
	"fmt"
	"templateGo/internal/model"
	"templateGo/internal/repositories"
)

// Example of how to use the new analytics tables
// IMPORTANT: These are TWO DIFFERENT types of statistics!

// CourseAnalyticsData represents COURSE-WIDE statistics (general course metrics)
// This includes data about the course as a whole, aggregated across all students
type CourseAnalyticsData struct {
	TotalEnrolledStudents    int     `json:"total_enrolled_students"`
	TotalActiveStudents      int     `json:"total_active_students"`
	CourseAverageGrade       float64 `json:"course_average_grade"`   // Average across ALL students
	CourseCompletionRate     float64 `json:"course_completion_rate"` // Percentage of students who completed
	TotalAssignmentsCreated  int     `json:"total_assignments_created"`
	MostDifficultAssignment  string  `json:"most_difficult_assignment"`
	MostPopularResource      string  `json:"most_popular_resource"`
	LastCourseActivity       string  `json:"last_course_activity"`
	InstructorEngagementRate float64 `json:"instructor_engagement_rate"`
}

// UserCourseAnalyticsData represents INDIVIDUAL USER statistics within a specific course
// This includes data specific to ONE user's performance and activity in the course
type UserCourseAnalyticsData struct {
	UserPersonalGrade        float64  `json:"user_personal_grade"`        // This specific user's grade
	UserAssignmentsCompleted int      `json:"user_assignments_completed"` // How many assignments THIS user completed
	UserTotalAssignments     int      `json:"user_total_assignments"`     // Total assignments available to THIS user
	UserTimeSpentMinutes     int      `json:"user_time_spent_minutes"`    // Time THIS user spent in the course
	UserLastActivity         string   `json:"user_last_activity"`         // When THIS user was last active
	UserPerformanceTrend     string   `json:"user_performance_trend"`     // "improving", "stable", "declining" for THIS user
	UserRankInCourse         int      `json:"user_rank_in_course"`        // THIS user's rank compared to other students
	UserFavoriteResources    []string `json:"user_favorite_resources"`    // Resources THIS user accesses most
	UserStudyPattern         string   `json:"user_study_pattern"`         // "morning", "evening", "weekend" for THIS user
}

func ExampleUsage() {
	// Get database connection
	db := repositories.GetDB()

	// Example 1: Create COURSE-WIDE analytics (statistics about the course as a whole)
	courseAnalytics := CourseAnalyticsData{
		TotalEnrolledStudents:    150,
		TotalActiveStudents:      142,
		CourseAverageGrade:       85.5, // Average grade across ALL students in the course
		CourseCompletionRate:     78.3, // Percentage of students who completed the course
		TotalAssignmentsCreated:  12,
		MostDifficultAssignment:  "Final Project",
		MostPopularResource:      "Chapter 5 Video",
		LastCourseActivity:       "2025-06-25T10:30:00Z",
		InstructorEngagementRate: 92.1,
	}

	// Convert to JSON string
	courseStatsJSON, _ := json.Marshal(courseAnalytics)

	// Save course-wide statistics to database
	courseAnalyticsRecord := model.CourseAnalytics{
		CourseID:   1, // Replace with actual course ID
		Statistics: string(courseStatsJSON),
	}

	result := db.Create(&courseAnalyticsRecord)
	if result.Error != nil {
		fmt.Printf("Error creating course analytics: %v\n", result.Error)
	}

	// Example 2: Create INDIVIDUAL USER analytics for a specific user in the course
	userAnalytics := UserCourseAnalyticsData{
		UserPersonalGrade:        92.5,                   // THIS user's grade (not the course average)
		UserAssignmentsCompleted: 8,                      // How many assignments THIS user completed
		UserTotalAssignments:     10,                     // Total assignments available to THIS user
		UserTimeSpentMinutes:     450,                    // Time THIS specific user spent
		UserLastActivity:         "2025-06-25T09:15:00Z", // When THIS user was last active
		UserPerformanceTrend:     "improving",            // THIS user's trend
		UserRankInCourse:         15,                     // THIS user's rank among all students
		UserFavoriteResources:    []string{"Video Lecture 3", "PDF Chapter 2"},
		UserStudyPattern:         "evening", // THIS user's preferred study time
	}

	// Convert to JSON string
	userStatsJSON, _ := json.Marshal(userAnalytics)

	// Save individual user statistics to database
	userAnalyticsRecord := model.UserCourseAnalytics{
		UserID:     "user123", // Replace with actual user ID
		CourseID:   1,         // Replace with actual course ID
		Statistics: string(userStatsJSON),
	}

	result = db.Create(&userAnalyticsRecord)
	if result.Error != nil {
		fmt.Printf("Error creating user analytics: %v\n", result.Error)
	}

	// Example 3: Query different types of analytics

	// Get COURSE-WIDE analytics
	var courseStats model.CourseAnalytics
	db.Where("course_id = ?", 1).First(&courseStats)

	// Parse JSON back to struct
	var parsedCourseStats CourseAnalyticsData
	json.Unmarshal([]byte(courseStats.Statistics), &parsedCourseStats)
	fmt.Printf("Course has %d enrolled students with course average grade %.2f\n",
		parsedCourseStats.TotalEnrolledStudents, parsedCourseStats.CourseAverageGrade)

	// Get INDIVIDUAL USER analytics for a specific user and course
	var userStats model.UserCourseAnalytics
	db.Where("user_id = ? AND course_id = ?", "user123", 1).First(&userStats)

	// Parse JSON back to struct
	var parsedUserStats UserCourseAnalyticsData
	json.Unmarshal([]byte(userStats.Statistics), &parsedUserStats)
	fmt.Printf("User has personal grade %.2f and completed %d/%d assignments (rank #%d in course)\n",
		parsedUserStats.UserPersonalGrade, parsedUserStats.UserAssignmentsCompleted,
		parsedUserStats.UserTotalAssignments, parsedUserStats.UserRankInCourse)
}
