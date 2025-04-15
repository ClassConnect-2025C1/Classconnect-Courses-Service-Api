package repositories

import "templateGo/internals/models"

// CourseRepositoryInterface define los métodos que debe implementar cualquier repositorio de cursos
type CourseRepositoryInterface interface {
	Create(course *models.Course) error

	GetByID(id uint) (*models.Course, error)

	GetAll() ([]models.Course, error)

	Update(course *models.Course) error

	Delete(id uint) error

	GetAvailableCourses(userID uint) ([]models.Course, error)

	IsUserEnrolled(courseID, userID uint) (bool, error)

	EnrollUser(courseID, userID uint, email, name string) error

	UnenrollUser(courseID, userID uint) error

	GetCourseMembers(courseID uint) ([]map[string]interface{}, error)

	UpdateMemberRole(courseID uint, userEmail string, role string) error
}
