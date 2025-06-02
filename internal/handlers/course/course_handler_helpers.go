package course

import (
	"net/http"
	"strconv"
	"templateGo/internal/model"
	"templateGo/internal/utils"

	"github.com/gin-gonic/gin"
)

// Parameter extraction helpers
func (h *courseHandlerImpl) getCourseID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Course ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandlerImpl) getAssignmentID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("assignment_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Assignment ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandlerImpl) getSubmissionID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Submission ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandlerImpl) getModuleID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("module_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "Module ID must be a number")
		return 0, false
	}
	return uint(id), true
}

func (h *courseHandlerImpl) getUserID(c *gin.Context) (string, bool) {
	id := c.Param("user_id")
	if id == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "Invalid Parameter", "User ID must be provided")
		return "", false
	}
	return id, true
}

func (h *courseHandlerImpl) getUserIDFromToken(c *gin.Context) (string, bool) {
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

func (h *courseHandlerImpl) getUserEmailFromToken(c *gin.Context) (string, bool) {
	user_email, exists := c.Get("user_email")
	if !exists {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User Email not found in token")
		return "", false
	}
	email, ok := user_email.(string)
	if !ok {
		utils.NewErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Invalid User Email format")
		return "", false
	}
	return email, true
}

// Entity retrieval helpers
func (h *courseHandlerImpl) getCourseByID(c *gin.Context, courseID uint) (*model.Course, bool) {
	course, err := h.repo.GetByID(courseID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Course not found")
		return nil, false
	}
	return course, true
}

func (h *courseHandlerImpl) getAssignmentByID(c *gin.Context, assignmentID uint) (*model.Assignment, bool) {
	assignment, err := h.repo.GetAssignmentByID(assignmentID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Assignment not found")
		return nil, false
	}
	return assignment, true
}

// TODO: que hcemos con esto?
func (h *courseHandlerImpl) getSubmissionByID(c *gin.Context, submissionID uint) (*model.Submission, bool) {
	submission, err := h.repo.GetSubmission(submissionID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Submission not found")
		return nil, false
	}
	return submission, true
}

func (h *courseHandlerImpl) getModuleByID(c *gin.Context, moduleID uint) (*model.Module, bool) {
	module, err := h.repo.GetModuleByID(moduleID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Module not found")
		return nil, false
	}
	return module, true
}

func (h *courseHandlerImpl) getResourceByID(c *gin.Context, resourceID string) (*model.Resource, bool) {
	resource, err := h.repo.GetResourceByID(resourceID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusNotFound, "Not Found", "Resource not found")
		return nil, false
	}
	return resource, true
}

// Response formatting helpers
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
		"teachingAssistants":  course.TeachingAssistants,
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func (h *courseHandlerImpl) isCourseCreatorOrAssistant(c *gin.Context, courseId uint) bool {
	userEmail, ok := h.getUserEmailFromToken(c)
	if !ok {
		return false
	}
	course, ok := h.getCourseByID(c, courseId)
	if !ok {
		return false
	}
	return course.CreatedBy == userEmail || contains(course.TeachingAssistants, userEmail)
}
