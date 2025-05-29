package course

import (
	"github.com/gin-gonic/gin"
)

// CourseHandler defines all course-related endpoint handlers
type CourseHandler interface {
	// Course Management
	CreateCourse(c *gin.Context)
	GetAllCourses(c *gin.Context)
	GetCourseByID(c *gin.Context)
	UpdateCourse(c *gin.Context)
	DeleteCourse(c *gin.Context)
	GetAvailableCourses(c *gin.Context)

	// Enrollment Management
	EnrollUserInCourse(c *gin.Context)
	UnenrollUserFromCourse(c *gin.Context)
	GetEnrolledCourses(c *gin.Context)
	GetCourseMembers(c *gin.Context)

	// Course Feedback Management
	CreateCourseFeedback(c *gin.Context)
	GetCourseFeedbacks(c *gin.Context)
	GetAIFeedbackAnalysis(c *gin.Context)

	// User Feedback Management
	CreateUserFeedback(c *gin.Context)
	GetUserFeedbacks(c *gin.Context)

	// Assignment Management
	CreateAssignment(c *gin.Context)
	UpdateAssignment(c *gin.Context)
	DeleteAssignment(c *gin.Context)
	GetAssignmentsPreviews(c *gin.Context)
	GetAssignmentByID(c *gin.Context)

	// Submission Management
	PutSubmissionOfCurrentUser(c *gin.Context)
	DeleteSubmissionOfCurrentUser(c *gin.Context)
	GetSubmissionOfCurrentUser(c *gin.Context)
	GetSubmissionByUserID(c *gin.Context)
	GetSubmissions(c *gin.Context)
	GradeSubmission(c *gin.Context)

	// Course Approval
	ApproveCourses(c *gin.Context)
	GetApprovedCourses(c *gin.Context)
	GetApprovedUsersForCourse(c *gin.Context) // New method

	// Course Favorites
	ToggleFavoriteStatus(c *gin.Context)
}
