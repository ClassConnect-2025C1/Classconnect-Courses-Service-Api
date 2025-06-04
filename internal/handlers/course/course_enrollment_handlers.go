package course

import (
	"errors"
	"net/http"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// EnrollUserInCourse handles user enrollment in a course
func (h *courseHandlerImpl) EnrollUserInCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	if err := h.repo.EnrollUser(courseID, userID); err != nil {
		if errors.Is(err, utils.ErrUserAlreadyEnrolled) {
			utils.NewErrorResponse(c, http.StatusConflict, "Conflict", "User is already enrolled in this course")
		} else {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error enrolling user in course")
		}
		return
	}

	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}
	h.notification.SendNotificationEmail(userID, course.Title)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled"})

}

// UnenrollUserFromCourse handles user unenrollment from a course
func (h *courseHandlerImpl) UnenrollUserFromCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	if err := h.repo.UnenrollUser(courseID, userID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error unenrolling user from course")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unenrolled"})

}

// GetEnrolledCourses returns courses the user is enrolled in
func (h *courseHandlerImpl) GetEnrolledCourses(c *gin.Context) {
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	courses, favorites, err := h.repo.GetEnrolledCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving enrolled courses")
		return
	}

	// Create a response that includes both course info and favorite status
	response := make([]map[string]any, len(courses))
	for i, course := range courses {
		courseMap := formatCourseResponse(&course)
		courseMap["is_favorite"] = favorites[i]
		response[i] = courseMap
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetCourseMembers returns all users enrolled in a course
func (h *courseHandlerImpl) GetCourseMembers(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	members, err := h.repo.GetCourseMembers(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving course members")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": members})
}
