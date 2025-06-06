package services

import (
	"fmt"
	"net/http"
	"templateGo/internal/metrics"

	"templateGo/internal/handlers/ai"
	"templateGo/internal/handlers/course"
	"templateGo/internal/handlers/notification"
	"templateGo/internal/logger"
	middleware "templateGo/internal/middlewares"
	"templateGo/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura las rutas del servidor y retorna un http.Handler.
func SetupRoutes(ddLogger *logger.DatadogLogger, ddMetrics *metrics.DatadogMetricsClient) http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Add middleware to log requests with Gin
	r.Use(func(c *gin.Context) {
		// Process request
		c.Next()

		// After request is processed
		if ddLogger != nil {
			status := c.Writer.Status()
			path := c.Request.URL.Path
			method := c.Request.Method

			attributes := map[string]interface{}{
				"status":    status,
				"path":      path,
				"method":    method,
				"client_ip": c.ClientIP(),
			}

			if status >= 400 {
				ddLogger.Error(fmt.Sprintf("%s %s - %d", method, path, status), attributes, nil)
			} else {
				ddLogger.Info(fmt.Sprintf("%s %s - %d", method, path, status), attributes, nil)
			}
		}
	})

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		if ddLogger != nil {
			ddLogger.Info("Health check request", nil, []string{"endpoint:health"})
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		fmt.Println("Response: healthcheck running wild")
	})

	// Create handlers with logger and metrics
	courseRepo := repositories.NewCourseRepository()
	notificationClient := notification.NewNotificationClient(nil)
	aiAnalyzer := ai.NewGeminiAnalyzer()
	courseHandler := course.NewCourseHandler(courseRepo, notificationClient, aiAnalyzer, ddMetrics)

	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		// =============================================
		// Course Management (CRUD operations)
		// =============================================

		// Create a new course
		api.POST("/course", courseHandler.CreateCourse)

		// Get all available courses
		api.GET("/courses", courseHandler.GetAllCourses)

		// Get details for a specific course
		api.GET("/:course_id", courseHandler.GetCourseByID)

		// Update course information
		api.PATCH("/:course_id", courseHandler.UpdateCourse)

		// Delete a course
		api.DELETE("/:course_id", courseHandler.DeleteCourse)

		// Get list of course members
		api.GET("/:course_id/members", courseHandler.GetCourseMembers)

		// Get available courses that a user can enroll in
		api.GET("/available", courseHandler.GetAvailableCourses)

		// Mark/unmark a course as favorite
		api.PATCH("/:course_id/favorite/toggle", courseHandler.ToggleFavoriteStatus)

		// =============================================
		// Enrollment Management
		// =============================================

		// Enroll current user in a course
		api.POST("/:course_id/enroll", courseHandler.EnrollUserInCourse)

		// Unenroll current user from a course
		api.DELETE("/:course_id/enroll", courseHandler.UnenrollUserFromCourse)

		// Get courses the current user is enrolled in
		api.GET("/enrolled", courseHandler.GetEnrolledCourses)

		// =============================================
		// Course Approval System
		// =============================================

		// Approve a user for a specific course
		api.POST("/approve/:user_id/:course_id", courseHandler.ApproveCourses)

		// Get approved courses for the current user
		api.GET("/approved", courseHandler.GetApprovedCourses)

		// Get approved users for a specific course
		api.GET("/:course_id/approved-users", courseHandler.GetApprovedUsersForCourse)

		// =============================================
		// Course Feedback & Ratings
		// =============================================

		// Submit feedback for a course
		api.POST("/:course_id/feedback", courseHandler.CreateCourseFeedback)

		// Get all feedback for a course
		api.GET("/:course_id/feedbacks", courseHandler.GetCourseFeedbacks)

		// Get AI-generated analysis of course feedback
		api.GET("/:course_id/ai-feedback-analysis", courseHandler.GetAIFeedbackAnalysis)

		// =============================================
		// User Feedback & Ratings
		// =============================================

		// Add feedback for a user in a course
		api.POST("/:course_id/user/:user_id/feedback", courseHandler.CreateUserFeedback)

		// Get all feedback for a user
		api.GET("/user/:user_id/feedbacks", courseHandler.GetUserFeedbacks)

		// =============================================
		// Assignment Management
		// =============================================

		// Create a new assignment for a course
		api.POST("/:course_id/assignment", courseHandler.CreateAssignment)

		// Get preview of all assignments in a course
		api.GET("/:course_id/assignments", courseHandler.GetAssignmentsPreviews)

		// Get details of a specific assignment
		api.GET("/:course_id/assignment/:assignment_id", courseHandler.GetAssignmentByID)

		// Update an existing assignment
		api.PATCH("/:course_id/assignment/:assignment_id", courseHandler.UpdateAssignment)

		// Delete an assignment
		api.DELETE("/:course_id/assignment/:assignment_id", courseHandler.DeleteAssignment)

		// =============================================
		// Submission Management
		// =============================================

		// Submit or update current user's assignment submission
		api.PUT("/:course_id/assignment/:assignment_id/submission", courseHandler.PutSubmissionOfCurrentUser)

		// Get current user's submission for an assignment
		api.GET("/:course_id/assignment/:assignment_id/submission", courseHandler.GetSubmissionOfCurrentUser)

		// Get all submissions for an assignment
		api.GET("/:course_id/assignment/:assignment_id/submissions", courseHandler.GetSubmissions)

		// Grade and provide feedback on a submission
		api.PATCH("/:course_id/assignment/:assignment_id/submission/:submission_id", courseHandler.GradeSubmission)

		// Get AI generated grade and feedback for a submission
		api.GET("/:course_id/assignment/:assignment_id/submission/:submission_id/ai-grade", courseHandler.GetAIGeneratedGradeAndFeedback)

		// Delete current user's submission
		api.DELETE("/:course_id/assignment/:assignment_id/submission", courseHandler.DeleteSubmissionOfCurrentUser)

		// =============================================
		// Resources Management
		// =============================================

		// Create a module for resources in a course
		api.POST("/:course_id/resource/module", courseHandler.CreateModule)

		// Create a resource in a specific module
		api.POST("/:course_id/resource/module/:module_id", courseHandler.CreateResource)

		// Patch a module name
		api.PATCH("/:course_id/resource/module/:module_id", courseHandler.PatchModule)

		// Get all resources(modules) from a course
		api.GET("/:course_id/resources", courseHandler.GetResources)

		// Patch order of modules and resources inside a course
		api.PATCH("/:course_id/resources", courseHandler.PatchResources)

		// Delete a resource in a specific module
		api.DELETE("/:course_id/resource/module/:module_id/:resource_id", courseHandler.DeleteResource)

		// Delete a module and all its resources
		api.DELETE("/:course_id/resource/module/:module_id", courseHandler.DeleteModule)

		// =============================================
		// Statistics
		// =============================================

		// Get statistics for all courses of the teacher
		api.GET("/statistics", courseHandler.GetCoursesStatistics)

		// Get statistics for a user
		api.GET("/statistics/course/:course_id/user/:user_id", courseHandler.GetUserStatisticsForCourse)
	}

	return r
}
