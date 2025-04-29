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
		api.POST("/courses", courseHandler.CreateCourse)
		api.GET("/courses", courseHandler.GetAllCourses)

		// Actualmente solo devuelve todos los cursos, deberia devolver los cursos
		// disponibles para el usuario autenticado en base a los criterios de elegibilidad
		api.GET("/available", courseHandler.GetAvailableCourses)

		// Rutas específicas por ID de curso
		api.GET("/:course_id", courseHandler.GetCourseByID)
		api.PATCH("/:course_id", courseHandler.UpdateCourse)
		api.DELETE("/:course_id", courseHandler.DeleteCourse)

		// Rutas de inscripción
		api.POST("/:course_id/enroll", courseHandler.EnrollUserInCourse)
		api.DELETE("/:course_id/enroll", courseHandler.UnenrollUserFromCourse)

		// Rutas de miembros
		api.GET("/:course_id/members", courseHandler.GetCourseMembers)
		api.PATCH("/:course_id/members/:user_email", courseHandler.UpdateMemberRole)

		// Rutas de feedback de cursos
		api.POST("/:course_id/feedback", courseHandler.CreateCourseFeedback)
		api.GET("/:course_id/feedback", courseHandler.GetCourseFeedback)
	}

	return r
}
