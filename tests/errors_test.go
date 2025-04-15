package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"templateGo/internals/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("GET", "/test/path", nil)
	return ctx, w
}

func TestErrorVariables(t *testing.T) {
	// Test that predefined errors exist with correct messages
	assert.Equal(t, "user already enrolled in this course", utils.ErrUserAlreadyEnrolled.Error())
	assert.Equal(t, "user not enrolled in this course", utils.ErrUserNotEnrolled.Error())
	assert.Equal(t, "course not found", utils.ErrCourseNotFound.Error())
}

func TestNewErrorResponse_BadRequest(t *testing.T) {
	c, w := setupTestContext()

	utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", "Invalid input data")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "https://api.classconnect.edu/errors/bad-request", response.Type)
	assert.Equal(t, "Validation Error", response.Title)
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "Invalid input data", response.Detail)
	assert.Equal(t, "/test/path", response.Instance)
}

func TestNewErrorResponse_Unauthorized(t *testing.T) {
	c, w := setupTestContext()

	utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "https://api.classconnect.edu/errors/unauthorized", response.Type)
	assert.Equal(t, "Unauthorized", response.Title)
}

func TestNewErrorResponse_NotFound(t *testing.T) {
	c, w := setupTestContext()

	utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Resource not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "https://api.classconnect.edu/errors/not-found", response.Type)
}

func TestNewErrorResponse_Conflict(t *testing.T) {
	c, w := setupTestContext()

	utils.NewErrorResponse(c, http.StatusConflict, "Conflict", "Resource already exists")

	assert.Equal(t, http.StatusConflict, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Note: Based on your implementation, conflict uses the server-error type
	assert.Equal(t, "https://api.classconnect.edu/errors/server-error", response.Type)
}

func TestNewErrorResponse_GenericError(t *testing.T) {
	c, w := setupTestContext()

	// Using a status code that doesn't have a specific mapping
	utils.NewErrorResponse(c, http.StatusTeapot, "I'm a teapot", "Cannot brew coffee")

	assert.Equal(t, http.StatusTeapot, w.Code)

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Should use the generic error type for unmapped status codes
	assert.Equal(t, "https://api.classconnect.edu/errors/generic", response.Type)
}

func TestErrorResponse_AllFieldsPopulated(t *testing.T) {
	c, w := setupTestContext()

	utils.NewErrorResponse(c, http.StatusBadRequest, "Test Title", "Test Detail")

	var response utils.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Ensure all fields are populated
	assert.NotEmpty(t, response.Type)
	assert.NotEmpty(t, response.Title)
	assert.NotEmpty(t, response.Detail)
	assert.NotEmpty(t, response.Instance)
	assert.NotZero(t, response.Status)
}
