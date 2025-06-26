package queue

import (
	"fmt"
	"templateGo/internal/handlers/ai"
	"templateGo/internal/repositories"

	"github.com/google/uuid"
)

// StatisticsService manages the statistics calculation queue
type StatisticsService struct {
	taskQueue *TaskQueue
}

// NewStatisticsService creates a new statistics service
func NewStatisticsService(repo repositories.CourseRepository, aiAnalyzer ai.FeedbackAnalyzer) *StatisticsService {
	// Create task processor
	processor := NewStatisticsTaskProcessor(repo, aiAnalyzer)

	// Create task queue with 3 workers and buffer size of 100
	taskQueue := NewTaskQueue(3, 100, processor)

	return &StatisticsService{
		taskQueue: taskQueue,
	}
}

// Start starts the statistics service
func (ss *StatisticsService) Start() {
	ss.taskQueue.Start()
}

// Stop stops the statistics service
func (ss *StatisticsService) Stop() {
	ss.taskQueue.Stop()
}

// EnqueueCourseStatisticsCalculation enqueues a course statistics calculation task
func (ss *StatisticsService) EnqueueCourseStatisticsCalculation(courseID uint, userID, userEmail string) error {
	task := Task{
		ID:   fmt.Sprintf("course-stats-%d-%s", courseID, uuid.New().String()[:8]),
		Type: TaskTypeCourseStatistics,
		Data: CourseStatisticsTaskData{
			CourseID:  courseID,
			UserID:    userID,
			UserEmail: userEmail,
		},
		MaxRetries: 3,
	}

	return ss.taskQueue.EnqueueTask(task)
}

// EnqueueUserCourseStatisticsCalculation enqueues a user course statistics calculation task
func (ss *StatisticsService) EnqueueUserCourseStatisticsCalculation(courseID uint, userID, userEmail string) error {
	task := Task{
		ID:   fmt.Sprintf("user-stats-%d-%s-%s", courseID, userID, uuid.New().String()[:8]),
		Type: TaskTypeUserCourseStatistics,
		Data: UserCourseStatisticsTaskData{
			CourseID:  courseID,
			UserID:    userID,
			UserEmail: userEmail,
		},
		MaxRetries: 3,
	}

	return ss.taskQueue.EnqueueTask(task)
}

// GetQueueSize returns the current queue size
func (ss *StatisticsService) GetQueueSize() int {
	return ss.taskQueue.GetQueueSize()
}
