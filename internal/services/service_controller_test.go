package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	sql "templateGo/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var router http.Handler

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	// Set up test environment
	gin.SetMode(gin.ReleaseMode)

	// Try to load .env from different possible locations
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load("../.env")
		if err != nil {
			fmt.Println("No .env file found, using environment variables")
		}
	}

	// Connect to test database
	if err := sql.ConnectDB(); err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}

	// Set up router
	router = SetupRoutes()

	// Run tests
	exitCode := m.Run()

	// Clean up
	sql.CloseDB()
	os.Exit(exitCode)
}

// Helper function to make API requests and parse responses
func makeRequest(method, url string, body any, target any) *httptest.ResponseRecorder {
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
	router.ServeHTTP(w, req)

	if target != nil && w.Code == http.StatusOK || w.Code == http.StatusCreated {
		err = json.Unmarshal(w.Body.Bytes(), target)
		if err != nil {
			panic(err)
		}
	}

	return w
}

func TestHealthcheck(t *testing.T) {
	// Setup router
	mux := SetupRoutes()

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

// CREATE COURSE TESTS
func TestCreateCourse_Success(t *testing.T) {
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    30,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, response["data"])
	data := response["data"].(map[string]any)
	assert.Equal(t, "Test Course", data["title"])

	// delete course
	courseID := data["id"].(string)
	makeRequest("DELETE", "/"+courseID, nil, nil)
}

func TestCreateCourse_MissingTitle(t *testing.T) {
	payload := map[string]any{
		"description": "Missing Title Course",
		"created_by":  "test@example.com",
		"capacity":    30,
	}

	w := makeRequest("POST", "/courses", payload, nil)

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
	makeRequest("POST", "/courses", payload, &response)

	payload2 := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    25,
	}

	var response2 map[string]any
	makeRequest("POST", "/courses", payload2, &response2)

	// La bdd deberia tener 2 cursos
	var response3 map[string]any
	w := makeRequest("GET", "/courses", nil, &response3)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response3["data"])
	courses := response3["data"].([]any)
	assert.GreaterOrEqual(t, len(courses), 2)

	// delete courses
	courseID := response["data"].(map[string]any)["id"].(string)
	courseID2 := response2["data"].(map[string]any)["id"].(string)
	makeRequest("DELETE", "/"+courseID, nil, nil)
	makeRequest("DELETE", "/"+courseID2, nil, nil)
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
	makeRequest("POST", "/courses", payload, &response)

	data := response["data"].(map[string]any)

	createdCourseID := data["id"].(string)

	// Get all courses
	w := makeRequest("GET", "/courses", nil, &response)
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
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
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
	makeRequest("POST", "/courses", payload, &response)

	data := response["data"].(map[string]any)
	createdCourseID := data["id"].(string)
	// Get the course by ID
	w := makeRequest("GET", "/"+createdCourseID, nil, &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response["data"])
	course := response["data"].(map[string]any)
	assert.Equal(t, createdCourseID, course["id"])
	assert.Equal(t, "Test Course", course["title"])
	assert.Equal(t, "Test Description", course["description"])

	// delete course
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

func TestGetCourseById_NotFound(t *testing.T) {
	w := makeRequest("GET", "/999999", nil, nil)
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
	makeRequest("POST", "/courses", payload, &response)

	data := response["data"].(map[string]any)
	createdCourseID := data["id"].(string)

	// Update the course
	updatePayload := map[string]any{
		"title":       "Updated Course",
		"description": "Updated Description",
	}

	w := makeRequest("PATCH", "/"+createdCourseID, updatePayload, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var updatedResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID, nil, &updatedResponse)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, updatedResponse)
	assert.Equal(t, "Updated Course", updatedResponse["title"])
	assert.Equal(t, "Updated Description", updatedResponse["description"])
	assert.Equal(t, createdCourseID, updatedResponse["id"])

	// delete course
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

func TestUpdateCourse_NotFound(t *testing.T) {
	payload := map[string]any{
		"title": "This Will Fail",
	}

	w := makeRequest("PATCH", "/999999", payload, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAvailableCourses_Success(t *testing.T) {
	var response map[string]any
	w := makeRequest("GET", "/courses", nil, &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, response["data"])
}

func TestGetAvailableCourses_InvalidUserId(t *testing.T) {
	w := makeRequest("GET", "/available?user_id=invalid", nil, nil)
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
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	// inscribir a un alumno
	enrollmentPayload := map[string]any{
		"user_id": 1,
		"email":   "estudiante@universidad.edu",
		"name":    "Pedro",
	}

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// ver que el alumno esta en el curso
	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members := membersResponse["data"].([]any)

	assert.GreaterOrEqual(t, len(members), 1)
	assert.Equal(t, "estudiante@universidad.edu", members[0].(map[string]any)["email"])
	assert.Equal(t, "Pedro", members[0].(map[string]any)["name"])
	assert.Equal(t, "student", members[0].(map[string]any)["role"])

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

// func TestEnrollUserInCourse_AlreadyEnrolled(t *testing.T) {
// 	// crear un curso
// 	payload := map[string]any{
// 		"title":       "Test Course",
// 		"description": "Test Description",
// 		"created_by":  "test@example.com",
// 		"capacity":    10,
// 	}

// 	var response map[string]any
// 	w := makeRequest("POST", "/courses", payload, &response)
// 	assert.Equal(t, http.StatusCreated, w.Code)

// 	createdCourseID := response["data"].(map[string]any)["id"].(string)

// 	// inscribir a un alumno
// 	enrollmentPayload := map[string]any{
// 		"user_id": 1,
// 		"email":   "estudiante@universidad.edu",
// 		"name":    "Pedro",
// 	}

// 	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	// ver que el alumno esta en el curso
// 	var membersResponse map[string]any
// 	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.NotNil(t, membersResponse["data"])

// 	members := membersResponse["data"].([]any)

// 	assert.GreaterOrEqual(t, len(members), 1)

// 	// intentar inscribir al mismo alumno de nuevo
// 	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
// 	assert.Equal(t, http.StatusConflict, w.Code)

// 	// eliminar curso
// 	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
// }

// GET MEMBERS TESTS
func TestGetCourseMembers_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	// inscribir a un alumno
	enrollmentPayload := map[string]any{
		"user_id": 1,
		"email":   "estudiante@universidad.edu",
		"name":    "Pedro",
	}

	enrollmentPayload2 := map[string]any{
		"user_id": 2,
		"email":   "tengo@muchosueño.com",
		"name":    "Juan",
	}

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload2, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// ver que el alumno esta en el curso
	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotNil(t, membersResponse["data"])
	members := membersResponse["data"].([]any)
	assert.GreaterOrEqual(t, len(members), 2)

	// comparar los nombres de los inscriptos
	assert.Equal(t, "Pedro", members[0].(map[string]any)["name"])
	assert.Equal(t, "Juan", members[1].(map[string]any)["name"])

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

// UPDATE MEMBER ROLE TESTS
func TestUpdateMemberRole_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	// inscribir a un alumno
	enrollmentPayload := map[string]any{
		"user_id": 1,
		"email":   "estudiante@universidad.edu",
		"name":    "Pedro",
	}

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// ver el rol del alumno
	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members := membersResponse["data"].([]any)
	assert.Equal(t, "student", members[0].(map[string]any)["role"])

	// actualizar el rol del alumno
	updatePayload := map[string]any{
		"role": "teacher",
	}

	w = makeRequest("PATCH", "/"+createdCourseID+"/members/estudiante@universidad.edu", updatePayload, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NotNil(t, membersResponse["data"])

	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members = membersResponse["data"].([]any)
	assert.Equal(t, "teacher", members[0].(map[string]any)["role"])

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

func TestUpdateUnexistentMemberRole(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	// actualizar el rol del alumno inexistente
	updatePayload := map[string]any{
		"role": "teacher",
	}

	w = makeRequest("PATCH", "/"+createdCourseID+"/members/estudiante@universidad.edu", updatePayload, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

// UNENROLL TESTS
func TestUnenrollUserFromCourse_Success(t *testing.T) {
	// crear un curso
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    10,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	// inscribir a un alumno
	enrollmentPayload := map[string]any{
		"user_id": 1,
		"email":   "estudiante@universidad.edu",
		"name":    "Pedro",
	}

	w = makeRequest("POST", "/"+createdCourseID+"/enroll", enrollmentPayload, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// desinscribir al alumno
	unenrollPayload := map[string]any{
		"user_id": 1,
		"email":   "estudiante@universidad.edu",
		"name":    "Pedro",
	}

	// DELETE http://localhost:8080/courses20/enroll

	w = makeRequest("DELETE", "/"+createdCourseID+"/enroll", unenrollPayload, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var membersResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID+"/members", nil, &membersResponse)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, membersResponse["data"])

	members := membersResponse["data"].([]any)
	fmt.Println("Members after unenrollment:", members)
	assert.Equal(t, 0, len(members))

	// eliminar curso
	makeRequest("DELETE", "/"+createdCourseID, nil, nil)
}

// DELETE COURSE TESTS
func TestDeleteCourse_Success(t *testing.T) {
	payload := map[string]any{
		"title":       "Test Course",
		"description": "Test Description",
		"created_by":  "test@example.com",
		"capacity":    15,
	}

	var response map[string]any
	w := makeRequest("POST", "/courses", payload, &response)
	assert.Equal(t, http.StatusCreated, w.Code)

	createdCourseID := response["data"].(map[string]any)["id"].(string)

	assert.NotEmpty(t, createdCourseID)

	w = makeRequest("DELETE", "/"+createdCourseID, nil, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify course was deleted
	var deletedResponse map[string]any
	w = makeRequest("GET", "/"+createdCourseID, nil, &deletedResponse)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteCourse_NotFound(t *testing.T) {
	w := makeRequest("DELETE", "/999999", nil, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
