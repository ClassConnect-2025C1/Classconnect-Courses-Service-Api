package services

import (
	"templateGo/internals/models"
	"templateGo/internals/repositories"
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
