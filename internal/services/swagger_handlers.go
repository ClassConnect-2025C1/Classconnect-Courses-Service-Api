package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Health check endpoint
// @Description Get the health status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} model.HealthCheckResponse
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// SwaggerInfo holds exported Swagger Info
type SwaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerConfig returns the swagger configuration
func SwaggerConfig() *SwaggerInfo {
	return &SwaggerInfo{
		Version:     "1.0.0",
		Host:        "localhost:8080",
		BasePath:    "/",
		Schemes:     []string{"http", "https"},
		Title:       "ClassConnect Courses Service API",
		Description: "API for managing courses, enrollments, assignments, and more in ClassConnect platform",
	}
}
