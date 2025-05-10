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

	GetAssignments(courseID uint) ([]model.Assignment, error)

	ApproveCourse(userID string, courseID uint, courseName string) error

	GetApprovedCourses(userID string) ([]string, error)

	// Toggle favorite status (flip current value)
	ToggleFavoriteStatus(courseID uint, userID string) error

	PutSubmission(submission *model.Submission) error

	GetSubmissionByUserID(courseID, assignmentID uint, userID string) (*model.Submission, error)

	GetSubmission(submissionID uint) (*model.Submission, error)

	GetSubmissions(courseID, assignmentID uint) ([]model.Submission, error)

	DeleteSubmission(submissionID uint) error

	GetAssignmentByID(assignmentID uint) (*model.Assignment, error)
}
