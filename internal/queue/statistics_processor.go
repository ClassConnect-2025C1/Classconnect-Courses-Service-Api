package queue

import (
	"fmt"
	"log"
	"templateGo/internal/handlers/ai"
	"templateGo/internal/model"
	"templateGo/internal/repositories"

	"gonum.org/v1/gonum/stat"
)

// StatisticsTaskProcessor handles statistics calculation tasks
type StatisticsTaskProcessor struct {
	repo       repositories.CourseRepository
	aiAnalyzer ai.FeedbackAnalyzer
}

// NewStatisticsTaskProcessor creates a new statistics task processor
func NewStatisticsTaskProcessor(repo repositories.CourseRepository, aiAnalyzer ai.FeedbackAnalyzer) *StatisticsTaskProcessor {
	return &StatisticsTaskProcessor{
		repo:       repo,
		aiAnalyzer: aiAnalyzer,
	}
}

// ProcessTask processes a task based on its type
func (stp *StatisticsTaskProcessor) ProcessTask(task Task) error {
	switch task.Type {
	case TaskTypeCourseStatistics:
		return stp.processCourseStatisticsTask(task)
	case TaskTypeUserCourseStatistics:
		return stp.processUserCourseStatisticsTask(task)
	case TaskTypeGlobalStatistics:
		return stp.processGlobalStatisticsTask(task)
	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}
}

// processCourseStatisticsTask processes course statistics calculation
func (stp *StatisticsTaskProcessor) processCourseStatisticsTask(task Task) error {
	data, ok := task.Data.(CourseStatisticsTaskData)
	if !ok {
		return fmt.Errorf("invalid task data type for course statistics task")
	}

	log.Printf("Processing course statistics for course %d", data.CourseID)

	// This is the same logic as in the original CalculateAndStoreCourseStatistics function
	// but moved to the task processor
	return stp.calculateAndStoreCourseStatistics(data.CourseID, data.UserID, data.UserEmail)
}

// processUserCourseStatisticsTask processes user course statistics calculation
func (stp *StatisticsTaskProcessor) processUserCourseStatisticsTask(task Task) error {
	data, ok := task.Data.(UserCourseStatisticsTaskData)
	if !ok {
		return fmt.Errorf("invalid task data type for user course statistics task")
	}

	log.Printf("Processing user course statistics for user %s in course %d", data.UserID, data.CourseID)

	// This is the same logic as in the original CalculateAndStoreUserCourseStatistics function
	// but moved to the task processor
	return stp.calculateAndStoreUserCourseStatistics(data.CourseID, data.UserID, data.UserEmail)
}

// The following functions are copied from the course handler but adapted for the task processor
// They contain the same business logic but are now part of the background processing

func (stp *StatisticsTaskProcessor) calculateAndStoreCourseStatistics(courseID uint, userID string, userEmail string) error {
	course, err := stp.repo.GetByID(courseID)
	if err != nil {
		return fmt.Errorf("error retrieving course: %w", err)
	}
	studentsCount, err := stp.repo.GetStudentsCount(courseID)
	if err != nil {
		return fmt.Errorf("error retrieving students count: %w", err)
	}
	globalTotalAverageGrade := 0.0
	globalTotalSubmissionRate := 0.0
	globalAssignmentsWithGradesCount := 0.0
	statisticsForAssignments := make([]model.StatisticsForAssignment, 0)
	assignments, err := stp.repo.GetAssignmentsPreviews(course.ID, userID, userEmail)
	if err != nil && err.Error() != "record not found" {
		return fmt.Errorf("error retrieving assignments: %w", err)
	}
	for _, assignment := range assignments {
		submissions, err := stp.repo.GetSubmissions(course.ID, assignment.ID)
		if err != nil && err.Error() != "record not found" {
			return fmt.Errorf("error retrieving submissions: %w", err)
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
		stp.calculateTendencyAndAverageGrade(last10Statistics)

	suggestions, _ := stp.aiAnalyzer.GenerateCourseSuggestionsBasedOnStats(Last10AssignmentsAverageGradeTendency, Last10AssignmentsSubmissionRateTendency, Last10AssignmentsAverageGrade)

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

	err = stp.repo.SaveCourseStatistics(statistics, course.ID)
	if err != nil {
		return fmt.Errorf("error saving course statistics: %w", err)
	}

	log.Printf("Successfully calculated and stored course statistics for course %d", courseID)
	return nil
}

func (stp *StatisticsTaskProcessor) calculateAndStoreUserCourseStatistics(courseID uint, studentID string, userEmail string) error {
	totalGrades := 0.0
	totalSubmissionsCount := 0.0
	totalRatedSubmissionsCount := 0.0
	statisticsForDates := make([]model.StatisticsForAssignment, 0)
	assignments, err := stp.repo.GetAssignmentsPreviews(courseID, studentID, userEmail)
	if err != nil {
		return fmt.Errorf("error retrieving assignments: %w", err)
	}
	for _, assignment := range assignments {
		submission, err := stp.repo.GetSubmissionByUserID(courseID, assignment.ID, studentID)
		if err != nil && err.Error() != "record not found" {
			return fmt.Errorf("error retrieving submission: %w", err)
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
		stp.calculateTendencyAndAverageGrade(last10Statistics)

	userStatistics := model.UserCourseStatistics{
		AverageGrade:                            averageGrade,
		SubmissionRate:                          submissionRate,
		Last10AssignmentsAverageGradeTendency:   Last10AssignmentsAverageGradeTendency,
		Last10AssignmentsSubmissionRateTendency: Last10AssignmentsSubmissionRateTendency,
		StatisticsForAssignments:                statisticsForDates,
	}

	err = stp.repo.SaveUserCourseStatistics(userStatistics, courseID, studentID)
	if err != nil {
		return fmt.Errorf("error saving user course statistics: %w", err)
	}

	log.Printf("Successfully calculated and stored user course statistics for user %s in course %d", studentID, courseID)
	return nil
}

func (stp *StatisticsTaskProcessor) calculateTendencyAndAverageGrade(stats []model.StatisticsForAssignment) (string, string, float64) {
	n := len(stats)
	if n == 0 {
		return "stable", "stable", 0
	}

	graded_n := 0

	// calculate the average first
	averageGrade := 0.0
	for _, stat := range stats {
		if stat.AverageGrade > 0 {
			averageGrade += stat.AverageGrade
			graded_n += 1
		}
	}
	if graded_n > 0 {
		averageGrade /= float64(graded_n)
	}

	if n == 1 || graded_n <= 1 {
		return "stable", "stable", averageGrade
	}

	xGrade := make([]float64, graded_n)
	yGradeFiltered := make([]float64, graded_n)

	xSubmission := make([]float64, n)
	ySubmission := make([]float64, n)

	gradedIndex := 0
	for i := 0; i < n; i++ {
		if stats[i].AverageGrade > 0 {
			xGrade[gradedIndex] = float64(gradedIndex)
			yGradeFiltered[gradedIndex] = stats[i].AverageGrade
			gradedIndex++
		}
	}

	for i := 0; i < n; i++ {
		xSubmission[i] = float64(i)
		ySubmission[i] = stats[i].SubmissionRate
	}

	_, slopeGrade := stat.LinearRegression(xGrade, yGradeFiltered, nil, false)
	_, slopeSubmission := stat.LinearRegression(xSubmission, ySubmission, nil, false)

	classify := func(slope float64) string {
		// Handle NaN or infinite values
		if slope != slope || slope == 0 { // NaN check (NaN != NaN is true)
			return "stable"
		}

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

	gradesTendency := classify(slopeGrade)
	submissionTendency := classify(slopeSubmission)

	// Additional safety check - ensure we never return empty strings
	if gradesTendency == "" {
		gradesTendency = "stable"
	}
	if submissionTendency == "" {
		submissionTendency = "stable"
	}

	return gradesTendency, submissionTendency, averageGrade
}

// processGlobalStatisticsTask processes global statistics calculation
func (stp *StatisticsTaskProcessor) processGlobalStatisticsTask(task Task) error {
	data, ok := task.Data.(GlobalStatisticsTaskData)
	if !ok {
		return fmt.Errorf("invalid task data type for global statistics task")
	}

	log.Printf("Processing global statistics for teacher %s", data.TeacherEmail)

	return stp.calculateAndStoreGlobalStatistics(data.TeacherEmail)
}

// calculateAndStoreGlobalStatistics calculates and stores global statistics for a teacher
func (stp *StatisticsTaskProcessor) calculateAndStoreGlobalStatistics(teacherEmail string) error {
	// Get all courses for the teacher
	courses, err := stp.repo.GetCoursesForTeacher(teacherEmail)
	if err != nil {
		return fmt.Errorf("error retrieving courses for teacher %s: %w", teacherEmail, err)
	}

	if len(courses) == 0 {
		log.Printf("No courses found for teacher %s", teacherEmail)
		return nil
	}

	var totalAverageGrade float64
	var totalSubmissionRate float64
	validCourses := 0

	// Calculate averages across all courses
	for _, course := range courses {
		courseStats, err := stp.repo.GetCourseStatistics(course.ID)
		if err != nil {
			log.Printf("Skipping course %d, no statistics available: %v", course.ID, err)
			continue // Skip courses without statistics
		}

		totalAverageGrade += courseStats.GlobalAverageGrade
		totalSubmissionRate += courseStats.GlobalSubmissionRate
		validCourses++
	}

	if validCourses == 0 {
		log.Printf("No valid courses with statistics found for teacher %s", teacherEmail)
		return nil
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
	err = stp.repo.SaveGlobalStatistics(globalStats)
	if err != nil {
		return fmt.Errorf("error saving global statistics for teacher %s: %w", teacherEmail, err)
	}

	log.Printf("Successfully calculated and stored global statistics for teacher %s (%.2f%% avg grade, %.2f%% submission rate from %d courses)",
		teacherEmail, globalAverageGrade, globalSubmissionRate, validCourses)

	return nil
}
