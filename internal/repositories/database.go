package repositories

import (
	"templateGo/internal/model"

	"gorm.io/gorm"
)

// DB instance is the global database connection
var DB *gorm.DB

// DatabaseManager defines the interface for database operations
type DatabaseManager interface {
	// ConnectDB establishes connection to the database
	ConnectDB() error

	// CloseDB closes the database connection
	CloseDB() error
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return DB
}

// Models that will be auto-migrated
var ModelsToMigrate = []any{
	&model.Course{},
	&model.Enrollment{},
	&model.CourseFeedback{},
	&model.Assignment{},
	&model.CourseApproval{},
	&model.Submission{},
	&model.AssignmentSession{},
	&model.UserFeedback{},
	&model.Module{},
	&model.Resource{},
}
