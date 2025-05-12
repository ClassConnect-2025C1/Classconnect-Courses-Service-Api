package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"templateGo/internal/model"
	"templateGo/internal/repositories"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type courseHandler struct {
	repo repositories.CourseRepository
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// Add a struct for the request
type toggleFavoriteRequest struct {
	IsFavorite bool `json:"is_favorite"`
}

func NewCourseHandler(repo repositories.CourseRepository) *courseHandler {
	return &courseHandler{repo}
}

func (h *courseHandler) getCourseID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Course ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandler) getAssignmentID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("assignment_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandler) getSubmissionID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Submission ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandler) getUserID(c *gin.Context) (string, bool) {
	id := c.Param("user_id")
	if id == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID must be provided")
		return "", false
	}
	return id, true
}

func (h *courseHandler) getUserIDFromToken(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User ID not found in token")
		return "", false
	}
	id, ok := userID.(string)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid User ID format")
		return "", false
	}
	return id, true
}

func (h *courseHandler) getCourseByID(c *gin.Context, courseID uint) (*model.Course, bool) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return nil, false
	}
	return course, true
}

func (h *courseHandler) getAssignmentByID(c *gin.Context, assignmentID uint) (*model.Assignment, bool) {
	assignment, err := h.repo.GetAssignmentByID(assignmentID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return nil, false
	}
	return assignment, true
}

func (h *courseHandler) getSubmissionByID(c *gin.Context, submissionID uint) (*model.Submission, bool) {
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return nil, false
	}
	return submission, true
}

// formatCoursesResponse formats multiple courses for API response
func formatCoursesResponse(courses []model.Course) []gin.H {
	response := make([]gin.H, 0, len(courses))
	for _, course := range courses {
		response = append(response, formatCourseResponse(&course))
	}
	return response
}

func formatCourseResponse(course *model.Course) gin.H {
	return gin.H{
		"id":                  strconv.FormatUint(uint64(course.ID), 10),
		"title":               course.Title,
		"description":         course.Description,
		"createdBy":           course.CreatedBy,
		"capacity":            course.Capacity,
		"startDate":           course.StartDate.Format("2006-01-02"),
		"endDate":             course.EndDate.Format("2006-01-02"),
		"eligibilityCriteria": course.EligibilityCriteria,
	}
}

// Handler methods
func (h *courseHandler) CreateCourse(c *gin.Context) {
	var request model.CreateCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	course := request.ToModel()
	if err := h.repo.Create(course); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating course")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully", "id": formatCourseResponse(course)["id"]})
}

func (h *courseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.repo.GetAll()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(courses)})
}

func (h *courseHandler) GetCourseByID(c *gin.Context) {
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

func (h *courseHandler) UpdateCourse(c *gin.Context) {
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

func (h *courseHandler) DeleteCourse(c *gin.Context) {
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

func (h *courseHandler) GetAvailableCourses(c *gin.Context) {
	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	courses, err := h.repo.GetAvailableCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving available courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(courses)})
}

func (h *courseHandler) GetEnrolledCourses(c *gin.Context) {
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
	response := make([]map[string]interface{}, len(courses))
	for i, course := range courses {
		courseMap := formatCourseResponse(&course)
		courseMap["is_favorite"] = favorites[i]
		response[i] = courseMap
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *courseHandler) EnrollUserInCourse(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled"})
}

func (h *courseHandler) UnenrollUserFromCourse(c *gin.Context) {
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

func (h *courseHandler) GetCourseMembers(c *gin.Context) {
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

func (h *courseHandler) CreateCourseFeedback(c *gin.Context) {
	courseID, ok := h.getCourseID(c)

	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var req model.CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	userID, ok := h.getUserIDFromToken(c)
	if !ok {
		return
	}

	// Check if user is enrolled in the course
	isEnrolled, err := h.repo.IsUserEnrolled(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error checking enrollment status")
		return
	}

	if !isEnrolled {
		utils.NewErrorResponse(c, http.StatusForbidden, "Forbidden", "Only enrolled users can provide feedback")
		return
	}

	feedback := &model.CourseFeedback{
		CourseID: courseID,
		Rating:   req.Rating,
		Comment:  req.Comment,
		Summary:  req.Summary,
	}

	if err := h.repo.CreateFeedback(feedback); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating feedback")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Feedback submitted successfully"})
}

func (h *courseHandler) GetCourseFeedbacks(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	feedbackList, err := h.repo.GetFeedbacksForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving feedback")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": feedbackList})
}

func (h *courseHandler) CreateAssignment(c *gin.Context) {
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

func (h *courseHandler) UpdateAssignment(c *gin.Context) {
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

func (h *courseHandler) DeleteAssignment(c *gin.Context) {
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

func (h *courseHandler) GetAssignments(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	assignments, err := h.repo.GetAssignments(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving assignments")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignments})
}

// ApproveCourses approves a course for a user
func (h *courseHandler) ApproveCourses(c *gin.Context) {
	userID, ok := h.getUserIDFromToken(c)
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
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error approving course")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course approved successfully"})
}

// GetApprovedCourses gets all approved courses for a user
func (h *courseHandler) GetApprovedCourses(c *gin.Context) {
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

// ToggleFavoriteStatus handles toggling a course's favorite status for a user
func (h *courseHandler) ToggleFavoriteStatus(c *gin.Context) {
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

func (h *courseHandler) PutSubmissionOfCurrentUser(c *gin.Context) {
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
}

func (h *courseHandler) DeleteSubmissionOfCurrentUser(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"message": "Submission deleted successfully"})
}

func (h *courseHandler) GetSubmissionOfCurrentUser(c *gin.Context) {
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

func (h *courseHandler) GetSubmissionByUserID(c *gin.Context) {
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

func (h *courseHandler) GetSubmissions(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
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

func (h *courseHandler) GradeSubmission(c *gin.Context) {
	submissionID, ok := h.getSubmissionID(c)
	if !ok {
		return
	}
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return
	}
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
}
