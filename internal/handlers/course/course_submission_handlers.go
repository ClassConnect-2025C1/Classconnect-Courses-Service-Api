package course

import (
	"log"
	"net/http"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// PutSubmissionOfCurrentUser creates or updates a submission
// @Summary Submit or update current user's assignment submission
// @Description Submit or update the current user's submission for an assignment
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Param submission body model.SubmissionRequest true "Submission content"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submission [put]
func (h *courseHandlerImpl) PutSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}
	var req model.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	submission := &model.Submission{
		CourseID:     courseID,
		AssignmentID: assignmentID,
		UserID:       userID,
		Content:      req.Content,
		Files:        req.Files,
	}

	if err := h.repo.PutSubmission(submission); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating submission")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Submission created/updated successfully"})

	// Enqueue statistics calculation tasks
	h.statisticsService.EnqueueCourseStatisticsCalculation(courseID, userID, userEmail)
	h.statisticsService.EnqueueUserCourseStatisticsCalculation(courseID, userID, userEmail)
	h.enqueueGlobalStatisticsForAllTeachers(courseID)
}

// DeleteSubmissionOfCurrentUser removes a user's submission
// @Summary Delete current user's submission
// @Description Delete the current user's submission for an assignment
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Success 204 "Submission deleted successfully"
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submission [delete]
func (h *courseHandlerImpl) DeleteSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return
	}
	// Check if assignment exists
	assignment, ok := h.getAssignmentByID(c, assignmentID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return
	}
	// Check if assignment belongs to the course
	if assignment.CourseID != courseID {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment does not belong to this course")
		return
	}
	// First get the submission of the user to check if it exists
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	if err := h.repo.DeleteSubmission(submission.ID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error deleting submission")
		return
	}
	c.JSON(http.StatusNoContent, nil)

	// Enqueue statistics calculation tasks
	h.statisticsService.EnqueueCourseStatisticsCalculation(courseID, userID, userEmail)
	h.statisticsService.EnqueueUserCourseStatisticsCalculation(courseID, userID, userEmail)
	// Also enqueue global statistics calculation for all teachers (creator + teaching assistants)
	h.enqueueGlobalStatisticsForAllTeachers(courseID)
}

// GetSubmissionOfCurrentUser returns the current user's submission
// @Summary Get current user's submission for an assignment
// @Description Retrieve the current user's submission for a specific assignment
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submission [get]
func (h *courseHandlerImpl) GetSubmissionOfCurrentUser(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submission})
}

// GetSubmissionByUserID returns a specific user's submission
func (h *courseHandlerImpl) GetSubmissionByUserID(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmissionByUserID(courseID, assignmentID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submission})

}

// GetSubmissions returns all submissions for an assignment
// @Summary Get all submissions for an assignment
// @Description Retrieve all submissions for a specific assignment (teacher only)
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submissions [get]
func (h *courseHandlerImpl) GetSubmissions(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	if !h.isCourseCreatorOrAssistant(c, courseID) {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "You do not have permission to access this resource")
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	submissions, err := h.repo.GetSubmissions(courseID, assignmentID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving submissions")
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": submissions})
}

// GradeSubmission allows grading a submission with feedback
// @Summary Grade and provide feedback on a submission
// @Description Grade a student's submission and provide feedback
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Param submission_id path string true "Submission ID"
// @Param grade body model.GradeRequest true "Grade and feedback"
// @Success 204 "Submission graded successfully"
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submission/{submission_id} [patch]
func (h *courseHandlerImpl) GradeSubmission(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return
	}
	submissionID, ok := h.getSubmissionID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	studentID := submission.UserID
	var req model.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}
	submission.Grade = req.Grade
	submission.Feedback = req.Feedback

	if err := h.repo.PutSubmission(submission); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error grading submission")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission graded successfully"})

	// Enqueue statistics calculation tasks
	h.statisticsService.EnqueueCourseStatisticsCalculation(courseID, userID, userEmail)
	h.statisticsService.EnqueueUserCourseStatisticsCalculation(courseID, studentID, userEmail)
	// Also enqueue global statistics calculation for all teachers (creator + teaching assistants)
	h.enqueueGlobalStatisticsForAllTeachers(courseID)
}

// GetAIGeneratedGrade retrieves AI-generated grade for a submission
// @Summary Get AI generated grade and feedback for a submission
// @Description Get AI-powered grade and feedback suggestions for a submission
// @Tags submissions
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID"
// @Param assignment_id path string true "Assignment ID"
// @Param submission_id path string true "Submission ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Security BearerAuth
// @Router /{course_id}/assignment/{assignment_id}/submission/{submission_id}/ai-grade [get]
func (h *courseHandlerImpl) GetAIGeneratedGradeAndFeedback(c *gin.Context) {
	submissionID, ok := h.getSubmissionID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
	assignmentID, ok := h.getAssignmentID(c)
	if !ok {
		return
	}
	// Check if assignment exists
	assignment, ok := h.getAssignmentByID(c, assignmentID)
	if !ok {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return
	}
	// Simulate AI-generated grade retrieval
	aiGrade, aiFeedback, err := h.aiAnalyzer.GenerateGradeAndFeedback(assignment.Description, submission.Files)
	if err != nil {
		// print it to the console for debugging
		log.Printf("Error generating AI grade/feedback: %v", err)

		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error generating AI grade/feedback")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"grade":    aiGrade,
		"feedback": aiFeedback,
	}})
}
