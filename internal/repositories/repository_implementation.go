package repositories

import (
	"fmt"
	"templateGo/internal/model"
	"time"

	"templateGo/internal/utils"

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
	var tempCourses []model.Course
	if err := r.db.Table("courses").Where("deleted_at IS NULL").Find(&tempCourses).Error; err != nil {
		return nil, err
	}

	return tempCourses, nil
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
	// devolver todos los cursos en los que el alumno no esta inscripto
	if err := r.db.Where("id NOT IN (SELECT course_id FROM enrollments WHERE user_id = ?)", userID).
		Where("deleted_at IS NULL").Find(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}

// Obtener cursos en los que un usuario está inscrito
func (r *courseRepository) GetEnrolledCourses(userID string) ([]model.Course, []bool, error) {
	var enrollments []model.Enrollment
	if err := r.db.Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		return nil, nil, err
	}

	courseIDs := make([]uint, 0, len(enrollments))
	favoriteStatus := make(map[uint]bool)

	for _, enrollment := range enrollments {
		courseIDs = append(courseIDs, enrollment.CourseID)
		favoriteStatus[enrollment.CourseID] = enrollment.Favorite
	}

	if len(courseIDs) == 0 {
		return []model.Course{}, []bool{}, nil
	}

	var courses []model.Course
	if err := r.db.Where("id IN ?", courseIDs).Find(&courses).Error; err != nil {
		return nil, nil, err
	}

	// Create a slice of favorite status in the same order as courses
	favorites := make([]bool, len(courses))
	for i, course := range courses {
		favorites[i] = favoriteStatus[course.ID]
	}

	return courses, favorites, nil
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
func (r *courseRepository) EnrollUser(courseID uint, userID string) error {
	enrollment := model.Enrollment{
		CourseID: courseID,
		UserID:   userID,
	}

	if err := r.db.Where("course_id = ? AND user_id = ?", courseID, userID).First(&enrollment).Error; err == nil {
		return utils.ErrUserAlreadyEnrolled
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
			"user_id": e.UserID,
		})
	}

	return members, nil
}

func (r *courseRepository) CreateFeedback(feedback *model.CourseFeedback) error {
	return DB.Create(feedback).Error
}

func (r *courseRepository) GetFeedbacksForCourse(courseID uint) ([]model.CourseFeedback, error) {
	var feedback []model.CourseFeedback
	err := DB.Where("course_id = ?", courseID).Find(&feedback).Error
	return feedback, err
}

func (r *courseRepository) CreateAssignment(assignment *model.Assignment) error {
	return DB.Create(assignment).Error
}

func (r *courseRepository) UpdateAssignment(assignment *model.Assignment) error {
	var existingAssignment model.Assignment
	if err := DB.First(&existingAssignment, assignment.ID).Error; err != nil {
		return err
	}
	// Start a transaction for atomicity
	tx := DB.Begin()

	// Clear existing file associations
	if err := tx.Model(&existingAssignment).Association("Files").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	assignment.CreatedAt = existingAssignment.CreatedAt // Preserve created time

	if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(assignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *courseRepository) DeleteAssignment(assignmentID uint) error {
	return DB.Delete(&model.Assignment{}, assignmentID).Error
}

func (r *courseRepository) GetAssignmentsPreviews(courseID uint, userID string) ([]model.AssignmentPreview, error) {
	var assignments []model.Assignment
	err := DB.Where("course_id = ?", courseID).Preload("Files").Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	previews := make([]model.AssignmentPreview, len(assignments))
	course := model.Course{}
	err = DB.Where("id = ?", courseID).First(&course).Error
	for i, assignment := range assignments {
		// Get status: if there exists a submission of user for this assignment, then it is submitted
		// if it's not submitted but a session exists, then it is started
		// if there is no session, then it is pending
		// status is none if the userID is of the course creator
		var status string
		if course.CreatedBy == userID {
			status = "none"
		} else {
			var submissionCount int64
			DB.Model(&model.Submission{}).
				Where("course_id = ? AND assignment_id = ? AND user_id = ?", courseID, assignment.ID, userID).
				Count(&submissionCount)

			if submissionCount > 0 {
				status = "submitted"
			} else {
				var sessionCount int64
				DB.Model(&model.AssignmentSession{}).
					Where("assignment_id = ? AND user_id = ?", assignment.ID, userID).
					Count(&sessionCount)

				if sessionCount > 0 {
					status = "started"
				} else {
					status = "pending"
				}
			}
		}
		previews[i] = model.AssignmentPreview{
			ID:        assignment.ID,
			Title:     assignment.Title,
			Deadline:  assignment.Deadline,
			TimeLimit: assignment.TimeLimit,
			CreatedAt: assignment.CreatedAt,
			DeletedAt: assignment.DeletedAt,
			Status:    status,
		}
	}
	return previews, err
}

func (r *courseRepository) GetAssignments(courseID uint) ([]model.Assignment, error) {
	var assignments []model.Assignment
	err := DB.Where("course_id = ?", courseID).Preload("Files").Find(&assignments).Error
	return assignments, err
}

func (r *courseRepository) GetOrCreateAssignmentSession(userID string, assignmentID uint) (*model.AssignmentSession, error) {
	var session model.AssignmentSession
	err := DB.Where("user_id = ? AND assignment_id = ?", userID, assignmentID).First(&session).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		session.UserID = userID
		session.AssignmentID = assignmentID
		session.StartedAt = time.Now()
		err = DB.Create(&session).Error
	}
	return &session, err
}

// ApproveCourse approves a course for a user
func (r *courseRepository) ApproveCourse(userID string, courseID uint, courseName string) error {
	approval := model.CourseApproval{
		UserID:     userID,
		CourseID:   courseID,
		CourseName: courseName,
	}

	// check if the user is enrolled in the course
	var enrollment model.Enrollment
	result := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user is not enrolled in this course: %w", result.Error)
		}
		return fmt.Errorf("error checking enrollment: %w", result.Error)
	}

	// Check if the approval already exists
	var count int64
	r.db.Model(&model.CourseApproval{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("user already approved")
	}

	return r.db.Create(&approval).Error
}

// GetApprovedCourses gets all approved course names for a user
func (r *courseRepository) GetApprovedCourses(userID string) ([]string, error) {
	var approvals []model.CourseApproval
	err := r.db.Where("user_id = ?", userID).Find(&approvals).Error
	if err != nil {
		return nil, err
	}

	courseNames := make([]string, 0, len(approvals))
	for _, approval := range approvals {
		courseNames = append(courseNames, approval.CourseName)
	}

	return courseNames, nil
}

// ToggleFavoriteStatus toggles the favorite status of a course for a user
func (r *courseRepository) ToggleFavoriteStatus(courseID uint, userID string) error {
	// Check if the user is enrolled in the course
	var enrollment model.Enrollment
	result := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment)

	if result.Error != nil {
		return fmt.Errorf("user is not enrolled in this course: %w", result.Error)
	}

	// Toggle favorite status (flip current value)
	enrollment.Favorite = !enrollment.Favorite
	return r.db.Save(&enrollment).Error
}

func (r *courseRepository) PutSubmission(submission *model.Submission) error {
	// Check if the submission already exists
	var existingSubmission model.Submission
	result := DB.Where("course_id = ? AND assignment_id = ? AND user_id = ?",
		submission.CourseID, submission.AssignmentID, submission.UserID).First(&existingSubmission)

	if result.Error == nil {
		// Submission exists, update it
		submission.ID = existingSubmission.ID // Preserve the original ID
		// If files are provided, clear existing associations first
		if len(submission.Files) > 0 {
			// Clear the files association to remove old files
			if err := DB.Model(&existingSubmission).Association("Files").Clear(); err != nil {
				return err
			}
		}
		return DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(submission).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// Submission doesn't exist, create it
		return DB.Create(submission).Error
	}

	return result.Error
}

func (r *courseRepository) GetSubmissionByUserID(courseID, assignmentID uint, userID string) (*model.Submission, error) {
	var submission model.Submission
	err := DB.Where("course_id = ? AND assignment_id = ? AND user_id = ?", courseID, assignmentID, userID).Preload("Files").First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *courseRepository) GetSubmission(submissionID uint) (*model.Submission, error) {
	var submission model.Submission
	err := DB.Where("id = ?", submissionID).First(&submission).Error
	if err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *courseRepository) GetSubmissions(courseID, assignmentID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	err := DB.Where("course_id = ? AND assignment_id = ?", courseID, assignmentID).Preload("Files").Find(&submissions).Error
	if err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *courseRepository) DeleteSubmission(submissionID uint) error {
	// First, find the submission
	var submission model.Submission
	if err := r.db.Where("id = ?",
		submissionID).First(&submission).Error; err != nil {
		return err
	}

	// Clear association with files (this removes entries from the junction table)
	if err := r.db.Model(&submission).Association("Files").Clear(); err != nil {
		return err
	}

	return r.db.Delete(&submission).Error
}

func (r *courseRepository) GetAssignmentByID(assignmentID uint) (*model.Assignment, error) {
	var assignment model.Assignment
	err := DB.Where("id = ?", assignmentID).Preload("Files").First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *courseRepository) GradeSubmission(submissionID uint, grade uint, feedback string) error {
	return DB.Model(&model.Submission{}).Where("id = ?", submissionID).Updates(model.Submission{
		Grade:    grade,
		Feedback: feedback,
	}).Error
}

// GetApprovedUsersForCourse retrieves all users approved for a specific course
func (r *courseRepository) GetApprovedUsersForCourse(courseID uint) ([]string, error) {
	var approvals []model.CourseApproval
	err := r.db.Where("course_id = ?", courseID).Find(&approvals).Error
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(approvals))
	for _, approval := range approvals {
		userIDs = append(userIDs, approval.UserID)
	}

	return userIDs, nil
}

// CreateUserFeedback adds feedback for a user in a course
func (r *courseRepository) CreateUserFeedback(feedback *model.UserFeedback) error {
	return r.db.Create(feedback).Error
}

// GetUserFeedbacks retrieves all feedback for a specific user
func (r *courseRepository) GetUserFeedbacks(userID string) ([]model.UserFeedback, error) {
	var feedbacks []model.UserFeedback

	// Simplify the query first to debug the issue
	err := r.db.Where("student_id = ?", userID).
		Order("created_at desc").
		Find(&feedbacks).Error

	return feedbacks, err
}

// CreateModule creates a new module for a course
func (r *courseRepository) CreateModule(module *model.Module) error {
	var lastModule model.Module
	err := r.db.Where("course_id = ?", module.CourseID).
		Order("\"order\" DESC").First(&lastModule).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error retrieving last module: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		module.Order = 0
	} else {
		module.Order = lastModule.Order + 1
	}
	return r.db.Create(module).Error
}

// CreateResource creates a new resource in a specific module
func (r *courseRepository) CreateResource(resource *model.Resource) error {
	var lastResource model.Resource
	err := r.db.Where("module_id = ?", resource.ModuleID).
		Order("\"order\" DESC").First(&lastResource).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error retrieving last resource: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		resource.Order = 0
	} else {
		resource.Order = lastResource.Order + 1
	}
	return r.db.Create(resource).Error
}

// GetModuleByID retrieves a module by its ID
func (r *courseRepository) GetModuleByID(moduleID uint) (*model.Module, error) {
	var module model.Module
	err := r.db.Where("id = ?", moduleID).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

// GetModulesByCourseID retrieves all modules for a specific course
func (r *courseRepository) GetModulesByCourseID(courseID uint) ([]model.Module, error) {
	var modules []model.Module
	err := r.db.Where("course_id = ?", courseID).Order("\"order\" ASC").Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// GetResourcesByModuleID retrieves all resources for a specific module
func (r *courseRepository) GetResourcesByModuleID(moduleID uint) ([]model.Resource, error) {
	var resources []model.Resource
	err := r.db.Where("module_id = ?", moduleID).Find(&resources).Error
	if err != nil {
		return nil, err
	}
	return resources, nil
}

// DeleteResource deletes a resource by its ID
func (r *courseRepository) DeleteResource(resourceID string) error {
	var resource model.Resource
	if err := r.db.Where("id = ?", resourceID).First(&resource).Error; err != nil {
		return fmt.Errorf("error retrieving resource: %w", err)
	}

	return r.db.Delete(&resource).Error
}

// DeleteModule deletes a module and all its resources
func (r *courseRepository) DeleteModule(moduleID uint) error {
	var module model.Module
	if err := r.db.Where("id = ?", moduleID).First(&module).Error; err != nil {
		return fmt.Errorf("error retrieving module: %w", err)
	}
	if err := r.db.Where("module_id = ?", moduleID).Delete(&model.Resource{}).Error; err != nil {
		return fmt.Errorf("error deleting resources for module: %w", err)
	}
	return r.db.Delete(&module).Error
}
