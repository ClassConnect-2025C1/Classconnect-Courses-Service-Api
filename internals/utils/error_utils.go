package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Errores personalizados
var (
	ErrUserAlreadyEnrolled = errors.New("user already enrolled in this course")
	ErrUserNotEnrolled     = errors.New("user not enrolled in this course")
	ErrCourseNotFound      = errors.New("course not found")
)

// ErrorResponse matches the OpenAPI error schema
type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

// NewErrorResponse creates a standard error response
func NewErrorResponse(c *gin.Context, status int, title string, detail string) {
	// Create path for the error instance
	instance := c.Request.URL.Path

	// Common error types based on status code
	errorType := "https://api.classconnect.edu/errors/generic"

	switch status {
	case http.StatusBadRequest: // 400
		errorType = "https://api.classconnect.edu/errors/bad-request"
	case http.StatusUnauthorized: // 401
		errorType = "https://api.classconnect.edu/errors/unauthorized"
	case http.StatusNotFound: // 404
		errorType = "https://api.classconnect.edu/errors/not-found"
	case http.StatusConflict: // 409
		errorType = "https://api.classconnect.edu/errors/server-error"
	}

	// RFC 7807 Format
	response := ErrorResponse{
		Type:     errorType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}

	c.JSON(status, response)
}
