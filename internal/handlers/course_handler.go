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

type enrollmentRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type unenrollRequest struct {
	UserID uint `json:"user_id" binding:"required"`
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

func (h *courseHandler) getUserID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandler) getCourseByID(c *gin.Context, courseID uint) (*model.Course, bool) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return nil, false
	}
	return course, true
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

	// quizas devolver solo un ok
	c.JSON(http.StatusCreated, gin.H{"data": formatCourseResponse(course)})
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
	userID, ok := h.getUserID(c)
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

func (h *courseHandler) EnrollUserInCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	var req enrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Request", "Invalid enrollment data: "+err.Error())
		return
	}

	if err := h.repo.EnrollUser(courseID, req.UserID, req.Email, req.Name); err != nil {
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

	// First try JSON body
	var req unenrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID is required")
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID must be a number")
			return
		}

		if err := h.repo.UnenrollUser(courseID, uint(userID)); err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error unenrolling user from course")
			return
		}
	} else {
		// JSON binding succeeded
		if err := h.repo.UnenrollUser(courseID, req.UserID); err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error unenrolling user from course")
			return
		}
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

func (h *courseHandler) UpdateMemberRole(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userEmail := c.Param("user_email")
	if userEmail == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User email is required")
		return
	}

	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	if err := h.repo.UpdateMemberRole(courseID, userEmail, req.Role); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error updating member role")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *courseHandler) IsUserEnrolled(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	isEnrolled, err := h.repo.IsUserEnrolled(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error checking user enrollment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_enrolled": isEnrolled})
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

	// Check if user is enrolled in the course
	isEnrolled, err := h.repo.IsUserEnrolled(courseID, req.UserID)
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
		UserID:   req.UserID,
		Rating:   req.Rating,
		Comment:  req.Comment,
		Summary:  req.Summary,
	}

	if err := h.repo.CreateFeedback(feedback); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating feedback")
		return
	}

	// Return response without including userID
	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback submitted successfully",
		"data": gin.H{
			"id":         feedback.ID,
			"course_id":  feedback.CourseID,
			"rating":     feedback.Rating,
			"comment":    feedback.Comment,
			"summary":    feedback.Summary,
			"created_at": feedback.CreatedAt,
		},
	})
}

func (h *courseHandler) GetCourseFeedback(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// Check if course exists
	_, ok = h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	feedbackList, err := h.repo.GetFeedbackForCourse(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving feedback")
		return
	}

	// Filter out userIDs from response
	responseData := make([]gin.H, 0, len(feedbackList))
	for _, item := range feedbackList {
		responseData = append(responseData, gin.H{
			"id":         item.ID,
			"course_id":  item.CourseID,
			"rating":     item.Rating,
			"comment":    item.Comment,
			"summary":    item.Summary,
			"created_at": item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}
