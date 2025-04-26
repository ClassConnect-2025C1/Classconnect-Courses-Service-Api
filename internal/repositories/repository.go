package repositories

import "templateGo/internal/model"

type CourseRepository interface {
	Create(course *model.Course) error

	GetByID(id uint) (*model.Course, error)

	GetAll() ([]model.Course, error)

	Update(course *model.Course) error

	Delete(id uint) error

	GetAvailableCourses(userID uint) ([]model.Course, error)

	IsUserEnrolled(courseID, userID uint) (bool, error)

	EnrollUser(courseID, userID uint, email, name string) error

	UnenrollUser(courseID, userID uint) error

	GetCourseMembers(courseID uint) ([]map[string]interface{}, error)

	UpdateMemberRole(courseID uint, userEmail string, role string) error

	CreateFeedback(feedback *model.CourseFeedback) error

	GetFeedbackForCourse(courseID uint) ([]model.CourseFeedback, error)
}
