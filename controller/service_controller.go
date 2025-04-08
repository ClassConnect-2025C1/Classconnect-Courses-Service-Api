package controller

import (
	"fmt"
	"net/http"
	"templateGo/internals/handlers"
	"templateGo/internals/repositories"
	"templateGo/internals/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura las rutas del servidor y retorna un http.Handler.
func SetupRoutes() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Health check endpoint
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		fmt.Println("Response: healthcheck")
	})

	// Create a new course handler
	courseRepo := repositories.NewCourseRepository()
	courseService := services.NewCourseService(courseRepo)
	courseHandler := handlers.NewCourseHandler(courseService)

	// Course endpoints
	courseGroup := r.Group("/courses")
	{
		courseGroup.POST("", courseHandler.CreateCourse)
		courseGroup.GET("", courseHandler.GetAllCourses)
		courseGroup.GET("/:id", courseHandler.GetCourseByID)
		courseGroup.PUT("/:id", courseHandler.UpdateCourse)
		courseGroup.DELETE("/:id", courseHandler.DeleteCourse)
	}

	return r
}
