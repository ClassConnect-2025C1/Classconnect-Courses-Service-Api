package repositories

import (
	"templateGo/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialize in-memory test database
func setupTestDB(t *testing.T) {
	var err error
	DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err, "Should connect to in-memory database")

	// Run migrations
	err = DB.AutoMigrate(&model.Course{})
	assert.NoError(t, err, "Should migrate tables")
}

func TestSQLDatabaseConnection(t *testing.T) {
	// Use in-memory database instead of PostgreSQL
	setupTestDB(t)

	// Ensure DB is initialized
	assert.NotNil(t, DB, "SQL DB instance should not be nil")

	// Test database connectivity
	var result int
	err := DB.Raw("SELECT 1").Scan(&result).Error
	assert.NoError(t, err, "Simple query should execute without error")
	assert.Equal(t, 1, result, "Query should return expected result")
}

func TestCloseDBFailure(t *testing.T) {
	DB = nil
	err := CloseDB()
	assert.NoError(t, err, "No error when closing nil connection")
}
