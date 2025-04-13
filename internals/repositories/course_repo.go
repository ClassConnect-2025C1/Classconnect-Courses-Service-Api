package repositories

import (
	"templateGo/dbConfig/sql"
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

// Obtener cursos disponibles para un usuario (a implementar según criterios de elegibilidad)
func (r *CourseRepository) GetAvailableCourses(userID uint) ([]models.Course, error) {
	var courses []models.Course
	// Aquí implementarías la lógica para obtener cursos disponibles según el usuario
	// Por ahora, simplemente retornamos todos los cursos no eliminados
	err := r.db.Where("deleted_at IS NULL").Find(&courses).Error
	return courses, err
}

// Verificar si un usuario ya está inscrito en un curso
func (r *CourseRepository) IsUserEnrolled(courseID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND user_id = ?", courseID, userID).
		Count(&count).Error

	return count > 0, err
}

// Inscribir a un usuario en un curso
func (r *CourseRepository) EnrollUser(courseID, userID uint, email, name string) error {
	enrollment := models.Enrollment{
		CourseID: courseID,
		UserID:   userID,
		Role:     "student", // Default role
		Email:    email,     // Provided email
		Name:     name,      // Provided name
	}

	return r.db.Create(&enrollment).Error
}

// Desinscribir a un usuario de un curso
func (r *CourseRepository) UnenrollUser(courseID, userID uint) error {
	return r.db.Where("course_id = ? AND user_id = ?", courseID, userID).
		Delete(&models.Enrollment{}).Error
}

// Obtener miembros de un curso
func (r *CourseRepository) GetCourseMembers(courseID uint) ([]map[string]interface{}, error) {
	var enrollments []models.Enrollment

	err := r.db.Where("course_id = ?", courseID).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	members := make([]map[string]interface{}, 0, len(enrollments))
	for _, e := range enrollments {
		members = append(members, map[string]any{
			"role":  e.Role,
			"name":  e.Name,
			"email": e.Email,
		})
	}

	return members, nil
}

// Actualizar rol de un miembro en un curso
func (r *CourseRepository) UpdateMemberRole(courseID uint, userEmail string, role string) error {
	return r.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND email = ?", courseID, userEmail).
		Update("role", role).Error
}
