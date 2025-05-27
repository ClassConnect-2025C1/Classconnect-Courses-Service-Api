package course

import (
	"net/http"
	"strconv"
	"strings"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// ApproveCourses approves a course for a user
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
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error approving course")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course approved successfully"})
}

// GetApprovedCourses returns all courses approved for a user
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
