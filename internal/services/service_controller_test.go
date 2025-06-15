package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"templateGo/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

var router http.Handler

// // TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	// Set up test environment
	gin.SetMode(gin.ReleaseMode)

	// Connect to test database
	dbManager := repositories.NewDatabaseManager()
	if err := dbManager.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up router
	router = SetupRoutes(nil, nil) // Pass nil for logger in tests

	// Run tests
	exitCode := m.Run()

	// Clean up
	defer dbManager.CloseDB()
	os.Exit(exitCode)
}

// // Helper function to make API requests and parse responses
func makeRequest(method, url string, body any, target any, user_id string, user_email string) *httptest.ResponseRecorder {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	token := generateTestToken(user_id, user_email)
	req.Header.Set("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)

	if target != nil && (w.Code == http.StatusOK || w.Code == http.StatusCreated) {
		err = json.Unmarshal(w.Body.Bytes(), target)
		if err != nil {
			panic(err)
		}
	}

	return w
}

func generateTestToken(user_id string, user_email string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["user_email"] = user_email

	secret_key := os.Getenv("JWT_SECRET_KEY")
	if secret_key == "" {
		secret_key = "supersecret"
	}

	tokenString, _ := token.SignedString([]byte(secret_key))

	return tokenString
}

func TestHealthcheck(t *testing.T) {
	// Setup router
	mux := SetupRoutes(nil, nil)

	// Create a test server (no need to specify port)
	server := httptest.NewServer(mux)
	defer server.Close() // This ensures the server is shut down after the test

	// Test "/healthcheck" endpoint using the server's URL
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Error making GET request to '/': %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for '/', got: %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got: %s", contentType)
	}

	// Check response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response from '/': %v", err)
	}
	expected := `{"status":"ok"}`
	if string(bodyBytes) != expected {
		t.Errorf("Response doesn't match, expected: %s, got: %s", expected, string(bodyBytes))
	}
}

// // CREATE COURSE TESTS
func TestCreateCourse_Success(t *testing.T) {
	payload := map[string]any{
		"title":                "Algo 4",
		"description":          "description test",
		"created_by":           "test01@gmail.com",
		"capacity":             10,
		"startDate":            "2025-06-01",
		"endDate":              "2025-09-30",
		"eligibility_criteria": []string{"Algo 3"},
		"teaching_assistants":  []string{"ta001", "ta002", "ta003"},
	}

	var response map[string]any
	w := makeRequest("POST", "/course", payload, &response, "user01", "test@example.com")

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, response["id"])

	courseID := response["id"].(string)

	// delete course
	makeRequest("DELETE", "/"+courseID, nil, nil, "", "")
}

func TestCreateCourse_MissingTitle(t *testing.T) {
	payload := map[string]any{
		"description": "Missing Title Course",
		"created_by":  "test@example.com",
		"capacity":    30,
	}

	w := makeRequest("POST", "/course", payload, nil, "user01", "test@example.com")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Validation Error")
}

func TestGetAllCourses(t *testing.T) {
	// Create a course to ensure there is at least one course in the database
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response map[string]any
	makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")

	payload2 := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response2 map[string]any
	makeRequest("POST", "/course", payload2, &response2, "user02", "test@example02.com")

	// La bdd deberia tener 2 cursos
	var response3 map[string]any
	w := makeRequest("GET", "/courses", nil, &response3, "", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response3["data"])
	courses := response3["data"].([]any)
	assert.GreaterOrEqual(t, len(courses), 2)

	// delete courses
	courseID := response["id"].(string)
	courseID2 := response2["id"].(string)
	makeRequest("DELETE", "/"+courseID, nil, nil, "", "")
	makeRequest("DELETE", "/"+courseID2, nil, nil, "", "")
}

func TestCreatedCourseExist(t *testing.T) {
	// Create a course to ensure there is at least one course in the database
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response map[string]any
	makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")

	createdCourseID := response["id"].(string)

	// Get all courses
	w := makeRequest("GET", "/courses", nil, &response, "user02", "test@example02.com")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response["data"])
	courses := response["data"].([]any)

	// Check if the created course is in the list
	found := false
	for _, course := range courses {
		c := course.(map[string]any)
		if c["id"] == createdCourseID {
			found = true
			break
		}
	}
	// delete course
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
	assert.True(t, found, "Created course should be in the list of all courses")
}

func TestGetCourseById_Success(t *testing.T) {
	// Create a course to ensure there is at least one course in the database
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response map[string]any
	makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")

	createdCourseID := response["id"].(string)

	// Get the course by ID
	var response2 map[string]any
	w := makeRequest("GET", "/"+createdCourseID, nil, &response2, "user01", "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response2)
	assert.Equal(t, createdCourseID, response2["id"])
	assert.Equal(t, "Test Course", response2["title"])
	assert.Equal(t, "Test Description", response2["description"])

	// delete course
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
}

func TestGetCourseById_NotFound(t *testing.T) {
	w := makeRequest("GET", "/999999", nil, nil, "", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateCourse_Success(t *testing.T) {
	// Create a course to ensure there is at least one course in the database
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response map[string]any
	makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")

	createdCourseID := response["id"].(string)

	// Update the course
	updatePayload := map[string]any{
		"title":       "Updated Course",
		"description": "Updated Description",
	}

	w := makeRequest("PATCH", "/"+createdCourseID, updatePayload, nil, "user01", "test@example01.com")
	assert.Equal(t, http.StatusNoContent, w.Code)

	var updatedResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID, nil, &updatedResponse, "user01", "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, updatedResponse)
	assert.Equal(t, "Updated Course", updatedResponse["title"])
	assert.Equal(t, "Updated Description", updatedResponse["description"])
	assert.Equal(t, createdCourseID, updatedResponse["id"])

	// // delete course
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
}

func TestUpdateCourse_NotFound(t *testing.T) {
	payload := map[string]any{
		"title": "This Will Fail",
	}

	w := makeRequest("PATCH", "/999999", payload, nil, "user01", "test@example01.com")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAvailableCourses_Success(t *testing.T) {
	var response map[string]any
	w := makeRequest("GET", "/courses", nil, &response, "user01", "test@example01.com")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response["data"])
}

func TestGetAvailableCourses_InvalidUserId(t *testing.T) {
	w := makeRequest("GET", "/available", nil, nil, "", "test@example01.com")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEnrollUserInCourse_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["id"].(string)
	userID := "1"

	// inscribir a un alumno usando la nueva ruta
	w = makeRequest("POST", "/"+createdCourseID+"/enroll", nil, nil, userID, "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	// ver que el alumno esta en el curso
	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse, "user01", "test@example01.com")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members := membersResponse["data"].([]any)

	assert.GreaterOrEqual(t, len(members), 1)
	assert.Equal(t, userID, members[0].(map[string]any)["user_id"])

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
}

func TestGetCourseMembers_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["id"].(string)

	// inscribir a los alumnos usando la nueva ruta
	userId1 := "1"
	userId2 := "2"

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", nil, nil, userId1, "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", nil, nil, userId2, "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	// ver que los alumnos están en el curso
	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse, "user01", "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotNil(t, membersResponse["data"])
	members := membersResponse["data"].([]any)
	assert.GreaterOrEqual(t, len(members), 2)

	// Recopilar los user_ids para verificación
	userIds := []string{}
	for _, member := range members {
		memberMap := member.(map[string]any)
		userIds = append(userIds, memberMap["user_id"].(string))
	}

	// Verificar que ambos user_ids estén presentes
	assert.Contains(t, userIds, userId1)
	assert.Contains(t, userIds, userId2)

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
}

func TestUnenrollUserFromCourse_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["id"].(string)
	userID := "1"

	// inscribir a un alumno usando la nueva ruta
	w = makeRequest("POST", "/"+createdCourseID+"/enroll", nil, nil, userID, "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	// desinscribir al alumno usando la nueva ruta
	w = makeRequest("DELETE", "/"+createdCourseID+"/enroll", nil, nil, userID, "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)

	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse, "user01", "test@example01.com")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members := membersResponse["data"].([]any)
	assert.Equal(t, 0, len(members))

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil, "", "")
}

func TestDeleteCourse_Success(t *testing.T) {
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    15,
	}

	var response map[string]any
	w := makeRequest("POST", "/course", payload, &response, "user01", "test@example01.com")
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["id"].(string)

	assert.NotEmpty(t, createdCourseID)

	w = makeRequest("DELETE", "/"+createdCourseID, nil, nil, "user01", "test@example01.com")
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify course was deleted
	var deletedResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID, nil, &deletedResponse, "user01", "test@example01.com")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteCourse_NotFound(t *testing.T) {
	w := makeRequest("DELETE", "/999999", nil, nil, "user01", "test@example01.com")
	assert.Equal(t, http.StatusNoContent, w.Code)
}
