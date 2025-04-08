package repositories

import (
	"templateGo/dbConfig/sql" // Update import to our new package
	"templateGo/internals/models"
	"time"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{db: sql.DB}
}

// Crear curso
func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// Obtener curso por ID (sin los eliminados)
func (r *CourseRepository) GetByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&course).Error
	return &course, err
}

// Obtener todos los cursos disponibles
func (r *CourseRepository) GetAll() ([]models.Course, error) {
	var courses []models.Course
	err := r.db.Where("deleted_at IS NULL").Find(&courses).Error
	return courses, err
}

// Editar curso
func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Eliminación lógica del curso
func (r *CourseRepository) Delete(id uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
