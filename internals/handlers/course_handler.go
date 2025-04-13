package handlers

import (
	"net/http"
	"strconv"
	"templateGo/internals/models"
	"templateGo/internals/services"
	"templateGo/internals/utils"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	service *services.CourseService
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type EnrollmentRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type UnenrollRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

func NewCourseHandler(service *services.CourseService) *CourseHandler {
	return &CourseHandler{service}
}

func (h *CourseHandler) getCourseID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Course ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *CourseHandler) getUserID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *CourseHandler) getCourseByID(c *gin.Context, courseID uint) (*models.Course, bool) {
	course, err := h.service.GetCourseByID(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return nil, false
	}
	return course, true
}

// formatCoursesResponse formats multiple courses for API response
func formatCoursesResponse(courses []models.Course) []gin.H {
	response := make([]gin.H, 0, len(courses))
	for _, course := range courses {
		response = append(response, formatCourseResponse(&course))
	}
	return response
}

func formatCourseResponse(course *models.Course) gin.H {
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
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var request models.CreateCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	course := request.ToModel()
	if err := h.service.CreateCourse(course); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error creating course")
		return
	}

	// quizas devolver solo un ok
	c.JSON(http.StatusCreated, gin.H{"data": formatCourseResponse(course)})
}

func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courses, err := h.service.GetAllCourses()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(courses)})
}

func (h *CourseHandler) GetCourseByID(c *gin.Context) {
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

func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	existingCourse, ok := h.getCourseByID(c, courseID)
	if !ok {
		return
	}

	var updateRequest models.UpdateCourseRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	updateRequest.ApplyTo(existingCourse)

	if err := h.service.UpdateCourse(existingCourse); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error updating course")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	if err := h.service.DeleteCourse(courseID); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error deleting course")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CourseHandler) GetAvailableCourses(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	courses, err := h.service.GetAvailableCourses(userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving available courses")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": formatCoursesResponse(courses)})
}

func (h *CourseHandler) EnrollUserInCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	var req EnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Request", "Invalid enrollment data: "+err.Error())
		return
	}

	if err := h.service.EnrollUser(courseID, req.UserID, req.Email, req.Name); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error enrolling user in course")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled"})
}

func (h *CourseHandler) UnenrollUserFromCourse(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	// First try JSON body
	var req UnenrollRequest
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

		if err := h.service.UnenrollUser(courseID, uint(userID)); err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error unenrolling user from course")
			return
		}
	} else {
		// JSON binding succeeded
		if err := h.service.UnenrollUser(courseID, req.UserID); err != nil {
			utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error unenrolling user from course")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unenrolled"})
}

func (h *CourseHandler) GetCourseMembers(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	members, err := h.service.GetCourseMembers(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error retrieving course members")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": members})
}

func (h *CourseHandler) UpdateMemberRole(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userEmail := c.Param("user_email")
	if userEmail == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User email is required")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	if err := h.service.UpdateMemberRole(courseID, userEmail, req.Role); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error updating member role")
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CourseHandler) IsUserEnrolled(c *gin.Context) {
	courseID, ok := h.getCourseID(c)
	if !ok {
		return
	}

	userID, ok := h.getUserID(c)
	if !ok {
		return
	}

	isEnrolled, err := h.service.IsUserEnrolled(courseID, userID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, "Server Error", "Error checking user enrollment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_enrolled": isEnrolled})
}
