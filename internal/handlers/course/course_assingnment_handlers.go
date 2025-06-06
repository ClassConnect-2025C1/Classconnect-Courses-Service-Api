package course

import (
	"net/http"
	"strconv"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// CreateAssignment creates a new assignment for a course
func (h *courseHandlerImpl) CreateAssignment(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var req model.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	assignment := &model.Assignment{
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		TimeLimit:   req.TimeLimit,
		Files:       req.Files,
	}

	if err := h.repo.CreateAssignment(assignment); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating assignment")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": assignment})
}

// UpdateAssignment updates an existing assignment
func (h *courseHandlerImpl) UpdateAssignment(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	assignmentID, err := strconv.Atoi(c.Param("assignment_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment ID must be a number")
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var req model.UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	assignment := &model.Assignment{
		ID:          uint(assignmentID),
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		TimeLimit:   req.TimeLimit,
		Files:       req.Files,
	}

	if err := h.repo.UpdateAssignment(assignment); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error updating assignment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignment})

}

// DeleteAssignment removes an assignment
func (h *courseHandlerImpl) DeleteAssignment(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	assignmentID, err := strconv.Atoi(c.Param("assignment_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment ID must be a number")
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	if err := h.repo.DeleteAssignment(uint(assignmentID)); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error deleting assignment")
		return
	}

	c.Status(http.StatusNoContent)

}

// GetAssignmentsPreviews returns previews of all assignments for a course
func (h *courseHandlerImpl) GetAssignmentsPreviews(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User email not found in context")
		return
	}

	assignments, err := h.repo.GetAssignmentsPreviews(courseID, userID, userEmail)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving assignments")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignments})

}

// GetAssignmentByID returns details of a specific assignment
func (h *courseHandlerImpl) GetAssignmentByID(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	assignmentID, err := strconv.Atoi(c.Param("assignment_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment ID must be a number")
		return
	}

	// Check if course exists
	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	// Get current user ID
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	assignment, ok := h.getAssignmentByID(c, uint(assignmentID))
	if !ok {
		return
	}

	// If the user is the teacher, just return the assignment
	if course.CreatedBy == userID {
		c.JSON(http.StatusOK, gin.H{"data": assignment})
		return
	}

	session, err := h.repo.GetOrCreateAssignmentSession(userID, uint(assignmentID))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error tracking assignment session")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": assignment,
		"session": gin.H{
			"started_at": session.StartedAt,
		},
	})

}
