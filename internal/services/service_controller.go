package services

import (
	"fmt"
	"net/http"
	"templateGo/internal/handlers"
	middleware "templateGo/internal/middlewares"
	"templateGo/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura las rutas del servidor y retorna un http.Handler.
func SetupRoutes() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		fmt.Println("Response: healthcheck running wild")
	})

	// Create a new course handler
	courseRepo := repositories.NewCourseRepository()
	courseHandler := handlers.NewCourseHandler(courseRepo)

	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/course", courseHandler.CreateCourse)
		api.GET("/courses", courseHandler.GetAllCourses)

		// Actualmente solo devuelve todos los cursos, deberia devolver los cursos
		// disponibles para el usuario autenticado en base a los criterios de elegibilidad
		api.GET("/available", courseHandler.GetAvailableCourses)
		api.GET("/enrolled", courseHandler.GetEnrolledCourses)

		// aprobar un usuario en un curso
		api.POST("/approve/:user_id/:course_id", courseHandler.ApproveCourses)

		// devolver los cursos que aprobo el usuario autenticado
		api.GET("/approved", courseHandler.GetApprovedCourses)

		// Rutas específicas por ID de curso
		api.GET("/:course_id", courseHandler.GetCourseByID)
		api.PATCH("/:course_id", courseHandler.UpdateCourse)
		api.DELETE("/:course_id", courseHandler.DeleteCourse)

		// Rutas de inscripción para usuario actual
		api.POST("/:course_id/enroll", courseHandler.EnrollUserInCourse)
		api.DELETE("/:course_id/enroll", courseHandler.UnenrollUserFromCourse)

		// Rutas de miembros
		api.GET("/:course_id/members", courseHandler.GetCourseMembers)

		// Rutas de feedback de cursos
		api.POST("/:course_id/feedback", courseHandler.CreateCourseFeedback)
		api.GET("/:course_id/feedbacks", courseHandler.GetCourseFeedbacks)

		// Add the new AI feedback analysis endpoint
		api.GET("/:course_id/ai-feedback-analysis", courseHandler.GetAIFeedbackAnalysis)

		// Rutas relacionadas a tareas(assignment)
		api.POST("/:course_id/assignment", courseHandler.CreateAssignment)
		api.PATCH("/:course_id/assignment/:assignment_id", courseHandler.UpdateAssignment)
		api.DELETE("/:course_id/assignment/:assignment_id", courseHandler.DeleteAssignment)
		api.GET("/:course_id/assignments", courseHandler.GetAssignmentsPreviews)
		api.GET("/:course_id/assignment/:assignment_id", courseHandler.GetAssignmentByID)

		// Rutas de submissions
		// Put/Delete submission of the current user
		api.PUT("/:course_id/assignment/:assignment_id/submission", courseHandler.PutSubmissionOfCurrentUser)
		api.DELETE("/:course_id/assignment/:assignment_id/submission", courseHandler.DeleteSubmissionOfCurrentUser)
		// Get submission of the current user
		api.GET("/:course_id/assignment/:assignment_id/submission", courseHandler.GetSubmissionOfCurrentUser)
		// Get submission of a specific user
		api.GET("/:course_id/assignment/:assignment_id/submission/:user_id", courseHandler.GetSubmissionByUserID)

		// Get all submissions
		api.GET("/:course_id/assignment/:assignment_id/submissions", courseHandler.GetSubmissions)
		// Grade and give feedback on a submission
		api.PATCH("/:course_id/assignment/:assignment_id/submission/:submission_id", courseHandler.GradeSubmission)

		// Toggle favorite status (switches between favorite and not favorite) of the current user
		api.PATCH("/:course_id/favorite/toggle", courseHandler.ToggleFavoriteStatus)
	}

	return r
}
