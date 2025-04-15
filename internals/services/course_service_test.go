package services

import (
	"errors"
	"testing"

	"templateGo/internals/models"
	"templateGo/internals/repositories"
	"templateGo/internals/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository implementation
type MockCourseRepository struct {
	mock.Mock
}

// Ensure MockCourseRepository implements CourseRepositoryInterface
var _ repositories.CourseRepositoryInterface = (*MockCourseRepository)(nil)

func (m *MockCourseRepository) Create(course *models.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByID(id uint) (*models.Course, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseRepository) GetAll() ([]models.Course, error) {
	args := m.Called()
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCourseRepository) Update(course *models.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCourseRepository) GetAvailableCourses(userID uint) ([]models.Course, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCourseRepository) IsUserEnrolled(courseID, userID uint) (bool, error) {
	args := m.Called(courseID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCourseRepository) EnrollUser(courseID, userID uint, email, name string) error {
	args := m.Called(courseID, userID, email, name)
	return args.Error(0)
}

func (m *MockCourseRepository) UnenrollUser(courseID, userID uint) error {
	args := m.Called(courseID, userID)
	return args.Error(0)
}

func (m *MockCourseRepository) GetCourseMembers(courseID uint) ([]map[string]interface{}, error) {
	args := m.Called(courseID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (m *MockCourseRepository) UpdateMemberRole(courseID uint, userEmail string, role string) error {
	args := m.Called(courseID, userEmail, role)
	return args.Error(0)
}

// Test cases for CourseService
func TestCreateCourse(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo)

	course := &models.Course{
		Title:       "Test Course",
		Description: "Test Description",
		CreatedBy:   "test@example.com",
		Capacity:    30,
	}

	mockRepo.On("Create", course).Return(nil)

	err := service.CreateCourse(course)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetCourseByID(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo)

	expectedCourse := &models.Course{
		Title:       "Test Course",
		Description: "Test Description",
		CreatedBy:   "test@example.com",
		Capacity:    30,
	}

	// Test success case
	mockRepo.On("GetByID", uint(1)).Return(expectedCourse, nil)
	course, err := service.GetCourseByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourse, course)

	// Test error case
	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))
	course, err = service.GetCourseByID(999)
	assert.Error(t, err)
	assert.Nil(t, course)

	mockRepo.AssertExpectations(t)
}

func TestEnrollUser(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo)

	// Test success case
	mockRepo.On("IsUserEnrolled", uint(1), uint(1)).Return(false, nil)
	mockRepo.On("EnrollUser", uint(1), uint(1), "test@example.com", "Test User").Return(nil)

	err := service.EnrollUser(1, 1, "test@example.com", "Test User")
	assert.NoError(t, err)

	// Test already enrolled case
	mockRepo.On("IsUserEnrolled", uint(1), uint(2)).Return(true, nil)

	err = service.EnrollUser(1, 2, "enrolled@example.com", "Enrolled User")
	assert.Equal(t, utils.ErrUserAlreadyEnrolled, err)

	mockRepo.AssertExpectations(t)
}

func TestUnenrollUser(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo)

	// Test success case
	mockRepo.On("IsUserEnrolled", uint(1), uint(1)).Return(true, nil)
	mockRepo.On("UnenrollUser", uint(1), uint(1)).Return(nil)

	err := service.UnenrollUser(1, 1)
	assert.NoError(t, err)

	// Test not enrolled case
	mockRepo.On("IsUserEnrolled", uint(1), uint(2)).Return(false, nil)

	err = service.UnenrollUser(1, 2)
	assert.NoError(t, err) // No error since trying to unenroll someone not enrolled is a no-op

	mockRepo.AssertExpectations(t)
}

func TestUpdateMemberRole(t *testing.T) {
	mockRepo := new(MockCourseRepository)
	service := NewCourseService(mockRepo)

	mockRepo.On("UpdateMemberRole", uint(1), "student@example.com", "teacher").Return(nil)

	err := service.UpdateMemberRole(1, "student@example.com", "teacher")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
