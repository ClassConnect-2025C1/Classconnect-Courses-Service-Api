package course

import (
	"net/http"
	"strconv"
	"strings"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// ApproveCourses approves a course for a user
// @Summary Approve a user for a specific course
// @Description Approve a user for enrollment in a specific course
// @Tags approval
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /approve/{user_id}/{course_id} [post]
func (h *courseHandlerImpl) ApproveCourses(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	// Get course_id from URL parameter
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 32)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Course ID must be a valid number")
		return
	}

	// Get the course to extract its name
	course, ok := h.getCourseByID(c, uint(courseID))
	if !ok {
		return
	}

	// Now use both the course ID and name
	if err := h.repo.ApproveCourse(userID, uint(courseID), course.Title); err != nil {
		if strings.Contains(err.Error(), "User already approved") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already approved"})
			return
		}

		if strings.Contains(err.Error(), "user is not enrolled in this course") {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not enrolled in this course"})
			return
		}

		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error approving course")
		return
	}

	h.notification.SendNotification(userID, course.Title, "course_approve")
	c.JSON(http.StatusOK, gin.H{"message": "Course approved successfully"})
}

// GetApprovedCourses returns all courses approved for a user
// @Summary Get approved courses for the current user
// @Description Retrieve all courses that the current user has been approved for
// @Tags approval
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /approved [get]
func (h *courseHandlerImpl) GetApprovedCourses(c *gin.Context) {
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	approvedCourses, err := h.repo.GetApprovedCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving approved courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": approvedCourses})
}

// ToggleFavoriteStatus toggles a course's favorite status for a user
// @Summary Mark/unmark a course as favorite
// @Description Toggle the favorite status of a course for the current user
// @Tags courses
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/favorite/toggle [patch]
func (h *courseHandlerImpl) ToggleFavoriteStatus(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	if err := h.repo.ToggleFavoriteStatus(courseID, userID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error toggling favorite status: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Course favorite status toggled successfully",
	})
}

// GetApprovedUsersForCourse returns all users approved for a specific course
// @Summary Get approved users for a specific course
// @Description Retrieve all users that have been approved for a specific course
// @Tags approval
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/approved-users [get]
func (h *courseHandlerImpl) GetApprovedUsersForCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Verify the course exists
	if _, ok := h.getCourseByID(c, courseID); !ok {
		return
	}

	approvedUsers, err := h.repo.GetApprovedUsersForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving approved users")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": approvedUsers})
}
