package course

import (
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
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
		globalAssignmentsWithGradesCount := 0.0
		statisticsForDates := make([]model.StatisticsForDate, 0)
		assignments, err := h.repo.GetAssignmentsPreviews(course.ID, userID, userEmail)
		if err != nil && err.Error() != "record not found" {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve assignments", "Error retrieving assignments: "+err.Error())
			return
		}
		for _, assignment := range assignments {
			submissions, err := h.repo.GetSubmissions(course.ID, assignment.ID)
			if err != nil && err.Error() != "record not found" {
				utils.NewErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve submissions", "Error retrieving submissions: "+err.Error())
				return
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
			averageGrade := 0.0
			submissionRate := 0.0
			if ratedSubmissionsCount > 0 {
				averageGrade = totalGrade / float64(ratedSubmissionsCount)
				globalAssignmentsWithGradesCount += 1
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
		if len(statisticsForDates) != 0 {
			globalTotalSubmissionRate /= float64(len(statisticsForDates))
		}
		if globalAssignmentsWithGradesCount != 0 {
			globalTotalAverageGrade /= globalAssignmentsWithGradesCount
		}

		last10Statistics := statisticsForDates
		if len(statisticsForDates) > 10 {
			last10Statistics = statisticsForDates[len(statisticsForDates)-10:]
		}
		Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency :=
			calculateTendency(last10Statistics)

		suggestions, _ := h.aiAnalyzer.GenerateCourseSuggestionsBasedOnStats(Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency)

		statistics = append(statistics, model.CourseStatistics{
			CourseID:                         course.ID,
			CourseName:                       course.Title,
			GlobalAverageGrade:               globalTotalAverageGrade,
			GlobalSubmissionRate:             globalTotalSubmissionRate,
			Last10DaysAverageGradeTendency:   Last10DaysAverageGradeTendency,
			Last10DaysSubmissionRateTendency: Last10DaysSubmissionRateTendency,
			Suggestions:                      suggestions,
			StatisticsForDates:               statisticsForDates,
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
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User email not found in context")
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
	assignments, err := h.repo.GetAssignmentsPreviews(courseID, userID, userEmail)
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
	averageGrade := 0.0
	submissionRate := 0.0
	assignmentsCount := len(assignments)
	if totalRatedSubmissionsCount > 0 {
		averageGrade = totalGrades / totalRatedSubmissionsCount
	}
	if assignmentsCount > 0 {
		submissionRate = totalSubmissionsCount / float64(assignmentsCount)
	}

	last10Statistics := statisticsForDates
	if len(statisticsForDates) > 10 {
		last10Statistics = statisticsForDates[len(statisticsForDates)-10:]
	}

	Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency :=
		calculateTendency(last10Statistics)

	userStatistics := model.UserCourseStatistics{
		AverageGrade:                     averageGrade,
		SubmissionRate:                   submissionRate,
		Last10DaysAverageGradeTendency:   Last10DaysAverageGradeTendency,
		Last10DaysSubmissionRateTendency: Last10DaysSubmissionRateTendency,
		StatisticsForDates:               statisticsForDates,
	}

	c.JSON(http.StatusOK, gin.H{"statistics": userStatistics})
}

func calculateTendency(stats []model.StatisticsForDate) (string, string) {
	n := len(stats)
	if n == 0 {
		return "stable", "stable"
	}

	x := make([]float64, n)
	yGrade := make([]float64, n)
	ySubmission := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = float64(i)
		yGrade[i] = stats[i].AverageGrade
		ySubmission[i] = stats[i].SubmissionRate
	}

	_, slopeGrade := stat.LinearRegression(x, yGrade, nil, false)
	_, slopeSubmission := stat.LinearRegression(x, ySubmission, nil, false)

	classify := func(slope float64) string {
		const epsilon = 0.01 // margen para considerar estable
		switch {
		case slope > epsilon:
			return "crecent"
		case slope < -epsilon:
			return "decrecent"
		default:
			return "stable"
		}
	}

	return classify(slopeGrade), classify(slopeSubmission)
}

func (h *courseHandlerImpl) CalculateAndStoreCourseStatistics(courseID uint, userID string, userEmail string) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		return
	}
	studentsCount, err := h.repo.GetStudentsCount(courseID)
	if err != nil {
		return
	}
	globalTotalAverageGrade := 0.0
	globalTotalSubmissionRate := 0.0
	globalAssignmentsWithGradesCount := 0.0
	statisticsForDates := make([]model.StatisticsForDate, 0)
	assignments, err := h.repo.GetAssignmentsPreviews(course.ID, userID, userEmail)
	for _, assignment := range assignments {
		submissions, err := h.repo.GetSubmissions(course.ID, assignment.ID)
		if err != nil && err.Error() != "record not found" {
			return
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
		averageGrade := 0.0
		submissionRate := 0.0
		if ratedSubmissionsCount > 0 {
			averageGrade = totalGrade / float64(ratedSubmissionsCount)
			globalAssignmentsWithGradesCount += 1
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
	if len(statisticsForDates) != 0 {
		globalTotalSubmissionRate /= float64(len(statisticsForDates))
	}
	if globalAssignmentsWithGradesCount != 0 {
		globalTotalAverageGrade /= globalAssignmentsWithGradesCount
	}

	last10Statistics := statisticsForDates
	if len(statisticsForDates) > 10 {
		last10Statistics = statisticsForDates[len(statisticsForDates)-10:]
	}
	Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency :=
		calculateTendency(last10Statistics)

	suggestions, _ := h.aiAnalyzer.GenerateCourseSuggestionsBasedOnStats(Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency)

	statistics := model.CourseStatistics{
		CourseID:                         course.ID,
		CourseName:                       course.Title,
		GlobalAverageGrade:               globalTotalAverageGrade,
		GlobalSubmissionRate:             globalTotalSubmissionRate,
		Last10DaysAverageGradeTendency:   Last10DaysAverageGradeTendency,
		Last10DaysSubmissionRateTendency: Last10DaysSubmissionRateTendency,
		Suggestions:                      suggestions,
		StatisticsForDates:               statisticsForDates,
	}

	err = h.repo.SaveCourseStatistics(statistics, course.ID)
	if err != nil {
		return
	}
}

func (h *courseHandlerImpl) CalculateAndStoreUserCourseStatistics(courseID uint, userID string, userEmail string) {
	totalGrades := 0.0
	totalSubmissionsCount := 0.0
	totalRatedSubmissionsCount := 0.0
	statisticsForDates := make([]model.StatisticsForDate, 0)
	assignments, err := h.repo.GetAssignmentsPreviews(courseID, userID, userEmail)
	if err != nil {
		return
	}
	for _, assignment := range assignments {
		submission, err := h.repo.GetSubmissionByUserID(courseID, assignment.ID, userID)
		if err != nil && err.Error() != "record not found" {
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
	averageGrade := 0.0
	submissionRate := 0.0
	assignmentsCount := len(assignments)
	if totalRatedSubmissionsCount > 0 {
		averageGrade = totalGrades / totalRatedSubmissionsCount
	}
	if assignmentsCount > 0 {
		submissionRate = totalSubmissionsCount / float64(assignmentsCount)
	}

	last10Statistics := statisticsForDates
	if len(statisticsForDates) > 10 {
		last10Statistics = statisticsForDates[len(statisticsForDates)-10:]
	}

	Last10DaysAverageGradeTendency, Last10DaysSubmissionRateTendency :=
		calculateTendency(last10Statistics)

	userStatistics := model.UserCourseStatistics{
		AverageGrade:                     averageGrade,
		SubmissionRate:                   submissionRate,
		Last10DaysAverageGradeTendency:   Last10DaysAverageGradeTendency,
		Last10DaysSubmissionRateTendency: Last10DaysSubmissionRateTendency,
		StatisticsForDates:               statisticsForDates,
	}

	err = h.repo.SaveUserCourseStatistics(userStatistics, courseID, userID)
	if err != nil {
		return
	}
}
