package services

import (
	"templateGo/internals/models"
	"templateGo/internals/repositories"
	"templateGo/internals/utils"
)

type CourseService struct {
	repo *repositories.CourseRepository
}

func NewCourseService(repo *repositories.CourseRepository) *CourseService {
	return &CourseService{repo}
}

// Crear curso
func (s *CourseService) CreateCourse(course *models.Course) error {
	return s.repo.Create(course)
}

// Obtener curso por ID
func (s *CourseService) GetCourseByID(id uint) (*models.Course, error) {
	return s.repo.GetByID(id)
}

// Obtener todos los cursos
func (s *CourseService) GetAllCourses() ([]models.Course, error) {
	return s.repo.GetAll()
}

// Editar curso
func (s *CourseService) UpdateCourse(course *models.Course) error {
	return s.repo.Update(course)
}

// Eliminar curso (lógicamente)
func (s *CourseService) DeleteCourse(id uint) error {
	return s.repo.Delete(id)
}

// Obtener cursos disponibles para un usuario
func (s *CourseService) GetAvailableCourses(userID uint) ([]models.Course, error) {
	return s.repo.GetAvailableCourses(userID)
}

// Inscribir a un usuario en un curso
func (s *CourseService) EnrollUser(courseID, userID uint, email, name string) error {
	// Verificamos si ya está inscrito
	enrolled, err := s.repo.IsUserEnrolled(courseID, userID)
	if err != nil {
		return err
	}
	if enrolled {
		return utils.ErrUserAlreadyEnrolled // Devuelve error específico
	}
	return s.repo.EnrollUser(courseID, userID, email, name)
}

// Desinscribir a un usuario de un curso
func (s *CourseService) UnenrollUser(courseID, userID uint) error {
	// Verificamos si está inscrito
	enrolled, err := s.repo.IsUserEnrolled(courseID, userID)
	if err != nil {
		return err
	}
	if !enrolled {
		return nil // Usuario no está inscrito, no hacemos nada
	}
	return s.repo.UnenrollUser(courseID, userID)
}

// Verificar si un usuario está inscrito en un curso
func (s *CourseService) IsUserEnrolled(courseID, userID uint) (bool, error) {
	return s.repo.IsUserEnrolled(courseID, userID)
}

// Obtener miembros de un curso
func (s *CourseService) GetCourseMembers(courseID uint) ([]map[string]any, error) {
	return s.repo.GetCourseMembers(courseID)
}

// Actualizar rol de un miembro en un curso
func (s *CourseService) UpdateMemberRole(courseID uint, userEmail string, role string) error {
	return s.repo.UpdateMemberRole(courseID, userEmail, role)
}
