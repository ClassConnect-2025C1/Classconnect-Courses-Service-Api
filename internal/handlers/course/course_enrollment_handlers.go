package course

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// EnrollUserInCourse handles user enrollment in a course
// @Summary Enroll the current user in a course
// @Description Enroll the authenticated user in the specified course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse{message=string}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/enroll [post]
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

	// Track enrollment metric
	if h.metricsClient != nil {
		// Add relevant tags for better filtering and visualization
		tags := []string{
			fmt.Sprintf("course_id:%d", courseID),
			fmt.Sprintf("course_title:%s", course.Title),
			fmt.Sprintf("user_id:%s", userID),
			fmt.Sprintf("environment:%s", os.Getenv("ENVIRONMENT")),
		}

		fmt.Printf("Sending metric: classconnect.course.enrollment.created with tags: %v\n", tags)
		if err := h.metricsClient.IncrementCounter("classconnect.course.enrollment.created", tags); err != nil {
			fmt.Printf("Error sending course enrollment metric: %v\n", err)
		} else {
			fmt.Println("Successfully sent enrollment metric to Datadog")
		}
	}

	//h.notification.SendNotificationEmail(userID, course.Title)
	h.notification.SendNotification(userID, course.Title, "enrollment")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled"})

}

// UnenrollUserFromCourse handles user unenrollment from a course
// @Summary Unenroll the current user from a course
// @Description Remove the authenticated user's enrollment from the specified course
// @Tags enrollments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse{message=string}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/enroll [delete]
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
// @Summary Get courses the current user is enrolled in
// @Description Retrieve all courses where the current user has an active enrollment
// @Tags enrollments
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse{data=[]model.CourseResponse}
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /enrolled [get]
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
// @Summary Retrieve members list for a course ID
// @Description Get all members enrolled in a specific course
// @Tags courses
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/members [get]
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
