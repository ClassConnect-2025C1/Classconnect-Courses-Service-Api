package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"templateGo/internal/externals"
	"templateGo/internal/model"
	"templateGo/internal/repositories"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type courseHandler struct {
	repo         repositories.CourseRepository
	notification *externals.NotificationClient
}

type updateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// Add a struct for the request
type toggleFavoriteRequest struct {
	IsFavorite bool `json:"is_favorite"`
}

func NewCourseHandler(repo repositories.CourseRepository, noti *externals.NotificationClient) *courseHandler {
	return &courseHandler{repo: repo, notification: noti}
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
	response := make([]map[string]any, len(courses))
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

	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}
	h.notification.SendNotificationEmail(userID, course.Title)
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

func (h *courseHandler) GetAssignmentsPreviews(c *gin.Context) {
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

	assignments, err := h.repo.GetAssignmentsPreviews(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving assignments")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": assignments})
}

func (h *courseHandler) GetAssignmentByID(c *gin.Context) {
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

// ApproveCourses approves a course for a user
func (h *courseHandler) ApproveCourses(c *gin.Context) {
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

// GetAIFeedbackAnalysis uses Gemini API to analyze course feedback
func (h *courseHandler) GetAIFeedbackAnalysis(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	course, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	// Get all feedback for this course
	feedbacks, err := h.repo.GetFeedbacksForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving feedback")
		return
	}

	// If there's no feedback, return an appropriate message
	if len(feedbacks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("No feedback available for analysis for course '%s'", course.Title),
		})
		return
	}

	// Format the feedback for the Gemini API
	feedbackText := formatFeedbackForAnalysis(course.Title, feedbacks)

	// Call Gemini API
	analysis, err := callGeminiAPI(feedbackText)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "AI Analysis Error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":           analysis,
		"feedback_count": len(feedbacks),
	})
}

// formatFeedbackForAnalysis formats the feedback data into text format for Gemini
func formatFeedbackForAnalysis(courseTitle string, feedbacks []model.CourseFeedback) string {
	var builder strings.Builder

	// builder.WriteString(fmt.Sprintf("Please analyze the following feedback for the course '%s' and provide a summary of common themes, strengths, and areas for improvement, make it short:\n\n", courseTitle))
	builder.WriteString(fmt.Sprintf("Please analyze the following feedback for the course '%s' and provide a summary of common themes, strengths, and areas for improvement, make it short, the rating is from 1 to 5, and I dont want any type of formatting in the text:", courseTitle))

	for i, feedback := range feedbacks {
		builder.WriteString(fmt.Sprintf("Feedback %d:\n", i+1))
		builder.WriteString(fmt.Sprintf("Rating: %d/100\n", feedback.Rating))
		if feedback.Summary != "" {
			builder.WriteString(fmt.Sprintf("Summary: %s\n", feedback.Summary))
		}
		if feedback.Comment != "" {
			builder.WriteString(fmt.Sprintf("Comment: %s\n", feedback.Comment))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// callGeminiAPI calls the Google Gemini API to get an analysis of the feedback
func callGeminiAPI(feedbackText string) (string, error) {
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	// apiKey := "AIzaSyCGf5mrU_9zlsOg538SsjJSeq1yIyyLXDc"

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	requestBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"text": feedbackText,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	// Extract the generated text from the response
	candidates, ok := responseData["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	candidate := candidates[0].(map[string]any)
	content, ok := candidate["content"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	parts, ok := content["parts"].([]any)
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	part := parts[0].(map[string]any)
	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	return text, nil
}
