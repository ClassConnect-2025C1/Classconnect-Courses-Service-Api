package services

import (
	"fmt"
	"net/http"
	"templateGo/internal/handlers"
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
	// api.Use(middleware.AuthMiddleware()) // Middleware for authentication if necessary
	{
		// Rutas según especificación OpenAPI
		api.POST("/course", courseHandler.CreateCourse)
		api.GET("/courses", courseHandler.GetAllCourses)

		// Actualmente solo devuelve todos los cursos, deberia devolver los cursos
		// disponibles para el usuario autenticado en base a los criterios de elegibilidad
		api.GET("/available/:user_id", courseHandler.GetAvailableCourses)
		api.GET("/enrolled/:user_id", courseHandler.GetEnrolledCourses)

		// aprobar un usuario en un curso
		api.GET("/approve/:user_id/:course_id", courseHandler.ApproveCourses)

		// devolver los cursos que aprobo el usuario autenticado
		api.GET("/approved/:user_id", courseHandler.GetApprovedCourses)

		// Rutas específicas por ID de curso
		api.GET("/:course_id", courseHandler.GetCourseByID)
		api.PATCH("/:course_id", courseHandler.UpdateCourse)
		api.DELETE("/:course_id", courseHandler.DeleteCourse)

		// Rutas de inscripción
		api.POST("/:course_id/enroll/:user_id", courseHandler.EnrollUserInCourse)
		api.DELETE("/:course_id/enroll/:user_id", courseHandler.UnenrollUserFromCourse)

		// Rutas de miembros
		api.GET("/:course_id/members", courseHandler.GetCourseMembers)

		// Rutas de feedback de cursos
		api.POST("/:course_id/feedback", courseHandler.CreateCourseFeedback)
		api.GET("/:course_id/feedbacks", courseHandler.GetCourseFeedbacks)

		// Rutas relacionadas a tareas(assignment)
		api.POST("/:course_id/assignment", courseHandler.CreateAssignment)
		api.PATCH("/:course_id/assignment/:assignment_id", courseHandler.UpdateAssignment)
		api.DELETE("/:course_id/assignment/:assignment_id", courseHandler.DeleteAssignment)
		api.GET("/:course_id/assignments", courseHandler.GetAssignments)

		// Rutas de submissions
		api.PUT("/:course_id/assignment/:assignment_id/submission", courseHandler.PutSubmission)
		api.DELETE("/:course_id/assignment/:assignment_id/submission/:user_id", courseHandler.DeleteSubmissionByUserID)
		// Get submission of a specific user
		api.GET("/:course_id/assignment/:assignment_id/submission/:user_id", courseHandler.GetSubmissionByUserID)

		// Get all submissions
		api.GET("/:course_id/assignment/:assignment_id/submissions", courseHandler.GetSubmissions)

		// Toggle favorite status (switches between favorite and not favorite)
		api.PATCH("/:course_id/favorite/toggle/:user_id", courseHandler.ToggleFavoriteStatus)
	}

	return r
}
