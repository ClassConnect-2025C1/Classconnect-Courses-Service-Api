package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	case http.StatusBadRequest:
		errorType = "https://api.classconnect.edu/errors/bad-request"
	case http.StatusUnauthorized:
		errorType = "https://api.classconnect.edu/errors/unauthorized"
	case http.StatusNotFound:
		errorType = "https://api.classconnect.edu/errors/not-found"
	case http.StatusInternalServerError:
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
