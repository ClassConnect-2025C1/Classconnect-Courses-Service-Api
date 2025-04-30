package repositories

import (
	"templateGo/internal/model"
	"time"

	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository() *courseRepository {
	return &courseRepository{db: DB}
}

// Crear curso
func (r *courseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

// Obtener curso por ID (sin los eliminados)
func (r *courseRepository) GetByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&course).Error
	return &course, err
}

// Obtener todos los cursos disponibles
func (r *courseRepository) GetAll() ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Where("deleted_at IS NULL").Find(&courses).Error
	return courses, err
}

// Editar curso
func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

// Eliminación lógica del curso
func (r *courseRepository) Delete(id uint) error {
	return r.db.Model(&model.Course{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

// Obtener cursos disponibles para un usuario (a implementar según criterios de elegibilidad)
func (r *courseRepository) GetAvailableCourses(userID string) ([]model.Course, error) {
	var courses []model.Course
	// Aquí implementarías la lógica para obtener cursos disponibles según el usuario
	// Por ahora, simplemente retornamos todos los cursos no eliminados
	err := r.db.Where("deleted_at IS NULL").Find(&courses).Error
	return courses, err
}

// Obtener cursos en los que un usuario está inscrito
func (r *courseRepository) GetEnrolledCourses(userID string) ([]model.Course, error) {
	var courses []model.Course

	var enrolledCourseIDs []uint
	err := r.db.Model(&model.Enrollment{}).
		Select("course_id").
		Where("user_id = ?", userID).
		Find(&enrolledCourseIDs).Error

	if err != nil {
		return nil, err
	}

	if len(enrolledCourseIDs) == 0 {
		return courses, nil
	}

	err = r.db.Where("id IN ? AND deleted_at IS NULL", enrolledCourseIDs).
		Find(&courses).Error

	return courses, err
}

// Verificar si un usuario ya está inscrito en un curso
func (r *courseRepository) IsUserEnrolled(courseID uint, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Enrollment{}).
		Where("course_id = ? AND user_id = ?", courseID, userID).
		Count(&count).Error

	return count > 0, err
}

// Inscribir a un usuario en un curso
func (r *courseRepository) EnrollUser(courseID uint, userID string, email, name string) error {
	enrollment := model.Enrollment{
		CourseID: courseID,
		UserID:   userID,
		Role:     "student", // Default role
		Email:    email,     // Provided email
		Name:     name,      // Provided name
	}

	return r.db.Create(&enrollment).Error
}

// Desinscribir a un usuario de un curso
func (r *courseRepository) UnenrollUser(courseID uint, userID string) error {
	return r.db.Where("course_id = ? AND user_id = ?", courseID, userID).
		Delete(&model.Enrollment{}).Error
}

// Obtener miembros de un curso
func (r *courseRepository) GetCourseMembers(courseID uint) ([]map[string]any, error) {
	var enrollments []model.Enrollment

	err := r.db.Where("course_id = ?", courseID).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}

	members := make([]map[string]any, 0, len(enrollments))
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
func (r *courseRepository) UpdateMemberRole(courseID uint, userEmail string, role string) error {
	return r.db.Model(&model.Enrollment{}).
		Where("course_id = ? AND email = ?", courseID, userEmail).
		Update("role", role).Error
}

func (r *courseRepository) CreateFeedback(feedback *model.CourseFeedback) error {
	return DB.Create(feedback).Error
}

func (r *courseRepository) GetFeedbackForCourse(courseID uint) ([]model.CourseFeedback, error) {
	var feedback []model.CourseFeedback
	err := DB.Where("course_id = ?", courseID).Find(&feedback).Error
	return feedback, err
}

func (r *courseRepository) CreateAssignment(assignment *model.Assignment) error {
	return DB.Create(assignment).Error
}

func (r *courseRepository) UpdateAssignment(assignment *model.Assignment) error {
	return DB.Save(assignment).Error
}

func (r *courseRepository) DeleteAssignment(assignmentID uint) error {
	return DB.Delete(&model.Assignment{}, assignmentID).Error
}

func (r *courseRepository) GetAssignments(courseID uint) ([]model.Assignment, error) {
	var assignments []model.Assignment
	err := DB.Where("course_id = ?", courseID).Find(&assignments).Error
	return assignments, err
}
