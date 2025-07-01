package course

import (
	"net/http"
	"strconv"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

const TIPO_NOTIFICACTION = "new_assignment"

// CreateAssignment creates a new assignment for a course
// @Summary Create a new assignment for a course
// @Description Create a new assignment within the specified course
// @Tags assignments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment body model.AssignmentRequest true "Assignment information"
// @Success 201 {object} model.SuccessResponse{data=model.AssignmentResponse}
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment [post]
func (h *courseHandlerImpl) CreateAssignment(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	course, ok := h.getCourseByID(c, courseID)
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

	courseMembers, _ := h.repo.GetCourseMembers(courseID)
	h.notification.SendNotificationToAll(courseMembers, course.Title, TIPO_NOTIFICACTION)
	c.JSON(http.StatusCreated, gin.H{"data": assignment})
}

// UpdateAssignment updates an existing assignment
// @Summary Update an existing assignment
// @Description Update the information of an existing assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Param assignment body model.AssignmentRequest true "Updated assignment information"
// @Success 204 "Assignment updated successfully"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id} [patch]
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
// @Summary Delete an assignment
// @Description Delete an assignment and all related submissions
// @Tags assignments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Success 204 "Assignment deleted successfully"
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id} [delete]
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
// @Summary Get preview of all assignments in a course
// @Description Retrieve a preview list of all assignments in a specific course
// @Tags assignments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignments [get]
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
// @Summary Get details of a specific assignment
// @Description Retrieve detailed information about a specific assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id} [get]
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
