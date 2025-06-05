package course

import (
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetCoursesStatistics retrieves statistics for all courses of the teacher (whether it's the creator or an teaching assistant)
func (h *courseHandlerImpl) GetCoursesStatistics(c *gin.Context) {
	// Get user email from the context
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User email not found in token")
		return
	}
	// Get courses for the teacher
	courses, err := h.repo.GetCoursesForTeacher(userEmail)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courses", "Error retrieving courses: "+err.Error())
		return
	}
	var statistics []model.CourseStatistics
	for _, course := range courses {
		studentsCount, err := h.repo.GetStudentsCount(course.ID)
		if err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve enrolled students count", "Error retrieving enrolled students count: "+err.Error())
			return
		}
		globalTotalAverageGrade := 0.0
		globalTotalSubmissionRate := 0.0
		statisticsForDates := make([]model.StatisticsForDate, 0)
		assignments, err := h.repo.GetAssignmentsPreviews(course.ID, userID)
		if err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve assignments", "Error retrieving assignments: "+err.Error())
			return
		}
		for _, assignment := range assignments {
			submissions, err := h.repo.GetSubmissions(course.ID, assignment.ID)
			if err != nil {
				utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve submissions", "Error retrieving submissions: "+err.Error())
				return
			}
			if len(submissions) == 0 {
				continue
			}
			totalGrade := 0.0
			submissionsCount := 0.0
			ratedSubmissionsCount := 0.0
			for _, submission := range submissions {
				submissionsCount += 1
				if submission.Grade > 0 {
					totalGrade += float64(submission.Grade)
					ratedSubmissionsCount += 1
				}
			}
			if submissionsCount == 0 {
				continue
			}
			averageGrade := 0.0
			submissionRate := 0.0
			if ratedSubmissionsCount > 0 {
				averageGrade = totalGrade / float64(ratedSubmissionsCount)
			}
			if studentsCount > 0 {
				submissionRate = submissionsCount / float64(studentsCount)
			}
			statisticsForDates = append(statisticsForDates, model.StatisticsForDate{
				Date:           assignment.CreatedAt,
				AverageGrade:   averageGrade,
				SubmissionRate: submissionRate,
			})
			globalTotalAverageGrade += averageGrade
			globalTotalSubmissionRate += submissionRate
		}
		if len(statisticsForDates) == 0 {
			continue
		}
		globalTotalAverageGrade /= float64(len(statisticsForDates))
		globalTotalSubmissionRate /= float64(len(statisticsForDates))
		statistics = append(statistics, model.CourseStatistics{
			CourseID:             course.ID,
			CourseName:           course.Title,
			GlobalAverageGrade:   globalTotalAverageGrade,
			GlobalSubmissionRate: globalTotalSubmissionRate,
			StatisticsForDates:   statisticsForDates,
		})
	}
	c.JSON(http.StatusOK, gin.H{"statistics": statistics})
}

// GetUserStatistics retrieves statistics for a specific user in a course
func (h *courseHandlerImpl) GetUserStatisticsForCourse(c *gin.Context) {
	// Get user email from the context
	userID, ok := h.getUserID(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in context")
		return
	}
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	if !h.isCourseCreatorOrAssistant(c, courseID) {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "You are not authorized to access this course statistics")
		return
	}

	totalGrades := 0.0
	totalSubmissionsCount := 0.0
	totalRatedSubmissionsCount := 0.0
	statisticsForDates := make([]model.StatisticsForDate, 0)
	assignments, err := h.repo.GetAssignmentsPreviews(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve assignments", "Error retrieving assignments: "+err.Error())
		return
	}
	for _, assignment := range assignments {
		submission, err := h.repo.GetSubmissionByUserID(courseID, assignment.ID, userID)
		if err != nil && err.Error() != "record not found" {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve submission", "Error retrieving submission: "+err.Error())
			return
		}
		if submission == nil {
			statisticsForDates = append(statisticsForDates, model.StatisticsForDate{
				Date:           assignment.CreatedAt,
				AverageGrade:   0.0,
				SubmissionRate: 0.0,
			})
			continue
		}
		totalSubmissionsCount += 1
		if submission.Grade > 0 {
			totalGrades += float64(submission.Grade)
			totalRatedSubmissionsCount += 1
		}
		statisticsForDates = append(statisticsForDates, model.StatisticsForDate{
			Date:           assignment.CreatedAt,
			AverageGrade:   float64(submission.Grade),
			SubmissionRate: 1.0,
		})
	}
	assignmentsCount := len(assignments)
	averageGrade := totalGrades / float64(totalRatedSubmissionsCount)
	submissionRate := totalSubmissionsCount / float64(assignmentsCount)

	userStatistics := model.UserCourseStatistics{
		AverageGrade:       averageGrade,
		SubmissionRate:     submissionRate,
		StatisticsForDates: statisticsForDates,
	}

	c.JSON(http.StatusOK, gin.H{"statistics": userStatistics})
}
