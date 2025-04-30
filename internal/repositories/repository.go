package repositories

import "templateGo/internal/model"

type CourseRepository interface {
	Create(course *model.Course) error

	GetByID(id uint) (*model.Course, error)

	GetAll() ([]model.Course, error)

	Update(course *model.Course) error

	Delete(id uint) error

	GetAvailableCourses(userID string) ([]model.Course, error)

	GetEnrolledCourses(userID string) ([]model.Course, error)

	IsUserEnrolled(courseID uint, userID string) (bool, error)

	EnrollUser(courseID uint, userID string, email, name string) error

	UnenrollUser(courseID uint, userID string) error

	GetCourseMembers(courseID uint) ([]map[string]interface{}, error)

	UpdateMemberRole(courseID uint, userEmail string, role string) error

	CreateFeedback(feedback *model.CourseFeedback) error

	GetFeedbacksForCourse(courseID uint) ([]model.CourseFeedback, error)

	CreateAssignment(assignment *model.Assignment) error

	UpdateAssignment(assignment *model.Assignment) error

	DeleteAssignment(assignmentID uint) error

	GetAssignments(courseID uint) ([]model.Assignment, error)
}
