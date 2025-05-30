package repositories

import "templateGo/internal/model"

type CourseRepository interface {
	Create(course *model.Course) error

	GetByID(id uint) (*model.Course, error)

	GetAll() ([]model.Course, error)

	Update(course *model.Course) error

	Delete(id uint) error

	GetAvailableCourses(userID string) ([]model.Course, error)

	GetEnrolledCourses(userID string) ([]model.Course, []bool, error)

	IsUserEnrolled(courseID uint, userID string) (bool, error)

	EnrollUser(courseID uint, userID string) error

	UnenrollUser(courseID uint, userID string) error

	GetCourseMembers(courseID uint) ([]map[string]any, error)

	CreateFeedback(feedback *model.CourseFeedback) error

	GetFeedbacksForCourse(courseID uint) ([]model.CourseFeedback, error)

	CreateAssignment(assignment *model.Assignment) error

	UpdateAssignment(assignment *model.Assignment) error

	DeleteAssignment(assignmentID uint) error

	GetAssignmentsPreviews(courseID uint, userID string) ([]model.AssignmentPreview, error)

	ApproveCourse(userID string, courseID uint, courseName string) error

	GetApprovedCourses(userID string) ([]string, error)

	ToggleFavoriteStatus(courseID uint, userID string) error

	PutSubmission(submission *model.Submission) error

	GetSubmissionByUserID(courseID, assignmentID uint, userID string) (*model.Submission, error)

	GetSubmission(submissionID uint) (*model.Submission, error)

	GetSubmissions(courseID, assignmentID uint) ([]model.Submission, error)

	DeleteSubmission(submissionID uint) error

	GetAssignmentByID(assignmentID uint) (*model.Assignment, error)

	GetOrCreateAssignmentSession(userID string, assignmentID uint) (*model.AssignmentSession, error)

	// GetApprovedUsersForCourse retrieves all users approved for a specific course
	GetApprovedUsersForCourse(courseID uint) ([]string, error)

	// CreateUserFeedback adds feedback for a user in a course
	CreateUserFeedback(feedback *model.UserFeedback) error

	// GetUserFeedbacks retrieves all feedback for a specific user
	GetUserFeedbacks(userID string) ([]model.UserFeedback, error)

	// CreateModule creates a new module for a course
	CreateModule(module *model.Module) error

	// CreateResource creates a new resource for a module
	CreateResource(resource *model.Resource) error

	GetModuleByID(moduleID uint) (*model.Module, error)
}
