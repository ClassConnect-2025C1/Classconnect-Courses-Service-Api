package course

import (
	"fmt"
	"net/http"
	"os"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateCourse handles course creation
func (h *courseHandlerImpl) CreateCourse(c *gin.Context) {
	var request model.CreateCourseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	fmt.Println("Creating course with request:", request)

	course := request.ToModel()
	if err := h.repo.Create(course); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating course")
		return
	}

	// Track course creation metric
	if h.metricsClient != nil {
		// Add relevant tags for better filtering and visualization
		tags := []string{
			fmt.Sprintf("course_id:%d", course.ID),
			fmt.Sprintf("course_name:%s", course.Title),
			fmt.Sprintf("course_type:%d", course.Capacity),
			fmt.Sprintf("environment:%s", os.Getenv("ENVIRONMENT")),
		}

		if err := h.metricsClient.IncrementCounter("classconnect.courses.created", tags); err != nil {
			fmt.Printf("Error sending course creation metric: %v\n", err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully", "id": formatCourseResponse(course)["id"]})

}

// GetAllCourses returns all courses
func (h *courseHandlerImpl) GetAllCourses(c *gin.Context) {
	courses, err := h.repo.GetAll()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(courses)})

}

// GetCourseByID returns a specific course by ID
func (h *courseHandlerImpl) GetCourseByID(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, formatCourseResponse(course))
}

// UpdateCourse updates an existing course
func (h *courseHandlerImpl) UpdateCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	existingCourse, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var updateRequest model.UpdateCourseRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	updateRequest.ApplyTo(existingCourse)

	if err := h.repo.Update(existingCourse); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error updating course")
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteCourse removes a course
func (h *courseHandlerImpl) DeleteCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	if err := h.repo.Delete(courseID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error deleting course")
		return
	}

	c.Status(http.StatusNoContent)

}

// GetAvailableCourses returns courses the user can enroll in based on eligibility criteria
func (h *courseHandlerImpl) GetAvailableCourses(c *gin.Context) {
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	// Get all courses the user is not enrolled in
	availableCourses, err := h.repo.GetAvailableCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving available courses")
		return
	}

	// Get the user's approved courses/subjects
	approvedSubjects, err := h.repo.GetApprovedCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving user's approved subjects")
		return
	}

	// Filter courses based on eligibility criteria
	var eligibleCourses []model.Course
	for _, course := range availableCourses {
		if meetsEligibilityCriteria(course, approvedSubjects) {
			eligibleCourses = append(eligibleCourses, course)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(eligibleCourses)})
}

// meetsEligibilityCriteria checks if a user with the given approved subjects meets a course's requirements
func meetsEligibilityCriteria(course model.Course, approvedSubjects []string) bool {
	// If there are no eligibility criteria, anyone can enroll
	if len(course.EligibilityCriteria) == 0 {
		return true
	}

	// Create a map for O(1) lookups of approved subjects
	approvedMap := make(map[string]bool)
	for _, subject := range approvedSubjects {
		approvedMap[subject] = true
	}

	// Check if all eligibility criteria are met
	for _, requirement := range course.EligibilityCriteria {
		if !approvedMap[requirement] {
			return false // User doesn't meet this requirement
		}
	}

	return true // All requirements are met
}
