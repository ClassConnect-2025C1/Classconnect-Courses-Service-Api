package course

import (
	"fmt"
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/stat"
)

// GetCoursesStatistics retrieves global statistics averages for the teacher
// @Summary Get global statistics averages for the teacher
// @Description Retrieve global average statistics across all courses taught by the current user
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /statistics/global [get]
func (h *courseHandlerImpl) GetCoursesStatistics(c *gin.Context) {
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User email not found in token")
		return
	}

	// Get global statistics for the teacher
	globalStats, err := h.repo.GetGlobalStatistics(userEmail)
	if err != nil {
		// Return default statistics if not found
		globalStats = model.GlobalStatistics{
			TeacherEmail:         userEmail,
			GlobalAverageGrade:   0.0,
			GlobalSubmissionRate: 0.0,
		}
	}

	c.JSON(http.StatusOK, gin.H{"statistics": globalStats})
}

// GetUserStatistics retrieves statistics for a specific user in a course
// @Summary Get statistics for a user in a course
// @Description Retrieve detailed statistics for a specific user within a course
// @Tags statistics
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /statistics/course/{course_id}/user/{user_id} [get]
func (h *courseHandlerImpl) GetUserStatisticsForCourse(c *gin.Context) {
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
	statistics, err := h.repo.GetUserCourseStatistics(courseID, userID)
	if err != nil {
		statistics = model.UserCourseStatistics{
			AverageGrade:                            0.0,
			SubmissionRate:                          0.0,
			Last10AssignmentsAverageGradeTendency:   "stable",
			Last10AssignmentsSubmissionRateTendency: "stable",
			StatisticsForAssignments:                []model.StatisticsForAssignment{},
		}
	}
	c.JSON(http.StatusOK, gin.H{"statistics": statistics})
}

// GetCourseStatistics retrieves statistics for a specific course
// @Summary Get statistics for a specific course
// @Description Retrieve comprehensive statistics for a specific course
// @Tags statistics
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /statistics/{course_id} [get]
func (h *courseHandlerImpl) GetCourseStatistics(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	if !h.isCourseCreatorOrAssistant(c, courseID) {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "You are not authorized to access this course statistics")
		return
	}

	courseStatistics, err := h.repo.GetCourseStatistics(courseID)
	if err != nil {
		// Get course info for the response even if no statistics are available
		course, courseErr := h.repo.GetByID(courseID)
		if courseErr != nil {
			utils.NewErrorResponse(c, http.StatusNotFound, "Course not found", "Course not found: "+courseErr.Error())
			return
		}

		courseStatistics = model.CourseStatistics{
			CourseID:                                courseID,
			CourseName:                              course.Title,
			GlobalAverageGrade:                      0.0,
			GlobalSubmissionRate:                    0.0,
			Suggestions:                             "No statistics available for this course",
			Last10AssignmentsAverageGradeTendency:   "stable",
			Last10AssignmentsSubmissionRateTendency: "stable",
			StatisticsForAssignments:                []model.StatisticsForAssignment{},
		}
	}

	c.JSON(http.StatusOK, gin.H{"statistics": courseStatistics})
}

func calculateTendencyAndAverageGrade(stats []model.StatisticsForAssignment) (string, string, float64) {
	n := len(stats)
	if n == 0 {
		return "stable", "stable", 0
	}
	// calculate the average first
	averageGrade := 0.0
	for _, stat := range stats {
		averageGrade += stat.AverageGrade
	}
	averageGrade /= float64(n)

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
			return "crescent"
		case slope < -epsilon:
			return "decrescent"
		default:
			return "stable"
		}
	}

	return classify(slopeGrade), classify(slopeSubmission), averageGrade
}

func (h *courseHandlerImpl) CalculateAndStoreCourseStatistics(courseID uint, userID string, userEmail string) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		fmt.Println("Error retrieving course:", err)
		return
	}
	studentsCount, err := h.repo.GetStudentsCount(courseID)
	if err != nil {
		fmt.Println("Error retrieving students count:", err)
		return
	}
	globalTotalAverageGrade := 0.0
	globalTotalSubmissionRate := 0.0
	globalAssignmentsWithGradesCount := 0.0
	statisticsForAssignments := make([]model.StatisticsForAssignment, 0)
	assignments, err := h.repo.GetAssignmentsPreviews(course.ID, userID, userEmail)
	if err != nil && err.Error() != "record not found" {
		fmt.Println("Error retrieving assignments:", err)
		return
	}
	for _, assignment := range assignments {
		submissions, err := h.repo.GetSubmissions(course.ID, assignment.ID)
		if err != nil && err.Error() != "record not found" {
			fmt.Println("Error retrieving submissions:", err)
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
		statisticsForAssignments = append(statisticsForAssignments, model.StatisticsForAssignment{
			Date:           assignment.CreatedAt,
			AverageGrade:   averageGrade,
			SubmissionRate: submissionRate,
		})
		globalTotalAverageGrade += averageGrade
		globalTotalSubmissionRate += submissionRate
	}
	if len(statisticsForAssignments) != 0 {
		globalTotalSubmissionRate /= float64(len(statisticsForAssignments))
	}
	if globalAssignmentsWithGradesCount != 0 {
		globalTotalAverageGrade /= globalAssignmentsWithGradesCount
	}

	last10Statistics := statisticsForAssignments
	if len(statisticsForAssignments) > 10 {
		last10Statistics = statisticsForAssignments[len(statisticsForAssignments)-10:]
	}
	Last10AssignmentsAverageGradeTendency, Last10AssignmentsSubmissionRateTendency, Last10AssignmentsAverageGrade :=
		calculateTendencyAndAverageGrade(last10Statistics)

	suggestions, _ := h.aiAnalyzer.GenerateCourseSuggestionsBasedOnStats(Last10AssignmentsAverageGradeTendency, Last10AssignmentsSubmissionRateTendency, Last10AssignmentsAverageGrade)

	statistics := model.CourseStatistics{
		CourseID:                                course.ID,
		CourseName:                              course.Title,
		GlobalAverageGrade:                      globalTotalAverageGrade,
		GlobalSubmissionRate:                    globalTotalSubmissionRate,
		Last10AssignmentsAverageGradeTendency:   Last10AssignmentsAverageGradeTendency,
		Last10AssignmentsSubmissionRateTendency: Last10AssignmentsSubmissionRateTendency,
		Suggestions:                             suggestions,
		StatisticsForAssignments:                statisticsForAssignments,
	}

	err = h.repo.SaveCourseStatistics(statistics, course.ID)
	if err != nil {
		fmt.Println("Error saving course statistics:", err)
		return
	}
}

func (h *courseHandlerImpl) CalculateAndStoreUserCourseStatistics(courseID uint, studentID string, userEmail string) {
	totalGrades := 0.0
	totalSubmissionsCount := 0.0
	totalRatedSubmissionsCount := 0.0
	statisticsForDates := make([]model.StatisticsForAssignment, 0)
	assignments, err := h.repo.GetAssignmentsPreviews(courseID, studentID, userEmail)
	if err != nil {
		return
	}
	for _, assignment := range assignments {
		submission, err := h.repo.GetSubmissionByUserID(courseID, assignment.ID, studentID)
		if err != nil && err.Error() != "record not found" {
			return
		}
		if submission == nil {
			statisticsForDates = append(statisticsForDates, model.StatisticsForAssignment{
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
		statisticsForDates = append(statisticsForDates, model.StatisticsForAssignment{
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

	Last10AssignmentsAverageGradeTendency, Last10AssignmentsSubmissionRateTendency, _ :=
		calculateTendencyAndAverageGrade(last10Statistics)

	userStatistics := model.UserCourseStatistics{
		AverageGrade:                            averageGrade,
		SubmissionRate:                          submissionRate,
		Last10AssignmentsAverageGradeTendency:   Last10AssignmentsAverageGradeTendency,
		Last10AssignmentsSubmissionRateTendency: Last10AssignmentsSubmissionRateTendency,
		StatisticsForAssignments:                statisticsForDates,
	}

	err = h.repo.SaveUserCourseStatistics(userStatistics, courseID, studentID)
	if err != nil {
		return
	}
}

// CalculateAndStoreGlobalStatistics calculates and stores global statistics for a teacher
func (h *courseHandlerImpl) CalculateAndStoreGlobalStatistics(teacherEmail string) {
	// Get all courses for the teacher
	courses, err := h.repo.GetCoursesForTeacher(teacherEmail)
	if err != nil {
		return
	}

	if len(courses) == 0 {
		return
	}

	var totalAverageGrade float64
	var totalSubmissionRate float64
	validCourses := 0

	// Calculate averages across all courses
	for _, course := range courses {
		courseStats, err := h.repo.GetCourseStatistics(course.ID)
		if err != nil {
			continue // Skip courses without statistics
		}

		totalAverageGrade += courseStats.GlobalAverageGrade
		totalSubmissionRate += courseStats.GlobalSubmissionRate
		validCourses++
	}

	if validCourses == 0 {
		return
	}

	// Calculate global averages (average of averages)
	globalAverageGrade := totalAverageGrade / float64(validCourses)
	globalSubmissionRate := totalSubmissionRate / float64(validCourses)

	// Create global statistics
	globalStats := model.GlobalStatistics{
		TeacherEmail:         teacherEmail,
		GlobalAverageGrade:   globalAverageGrade,
		GlobalSubmissionRate: globalSubmissionRate,
	}

	// Save global statistics
	err = h.repo.SaveGlobalStatistics(globalStats)
	if err != nil {
		return
	}
}

// getAllTeachersForCourse returns all teachers (creator + teaching assistants) for a course
func (h *courseHandlerImpl) getAllTeachersForCourse(courseID uint) ([]string, error) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		return nil, err
	}

	teachers := []string{course.CreatedBy}

	// Add teaching assistants
	for _, ta := range course.TeachingAssistants {
		teachers = append(teachers, ta)
	}

	return teachers, nil
}

// enqueueGlobalStatisticsForAllTeachers enqueues global statistics calculation for all teachers of a course
func (h *courseHandlerImpl) enqueueGlobalStatisticsForAllTeachers(courseID uint) {
	teachers, err := h.getAllTeachersForCourse(courseID)
	if err != nil {
		return
	}

	for _, teacherEmail := range teachers {
		h.statisticsService.EnqueueGlobalStatisticsCalculation(teacherEmail)
	}
}
