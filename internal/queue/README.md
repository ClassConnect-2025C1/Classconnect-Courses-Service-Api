# Statistics Queue System

This module implements a producer-consumer pattern queue system for processing statistics calculations asynchronously.

## Overview

The statistics queue system decouples HTTP request handling from potentially time-consuming statistics calculations by processing them in background worker goroutines.

## Components

### 1. Task Queue (`task_queue.go`)
- Generic task queue with configurable workers and buffer size
- Supports graceful shutdown and automatic retry logic
- Provides worker pool management and context-based cancellation

### 2. Task Processor (`statistics_processor.go`)
- Implements the business logic for calculating course and user statistics
- Processes tasks based on their type (course statistics or user statistics)
- Contains the same calculation logic as the original handlers but runs asynchronously

### 3. Statistics Service (`statistics_service.go`)
- High-level interface for enqueueing statistics calculation tasks
- Manages the task queue lifecycle (start/stop)
- Provides convenient methods for different types of statistics calculations

### 4. Task Data (`task_data.go`)
- Defines the data structures for different task types
- Contains CourseStatisticsTaskData and UserCourseStatisticsTaskData

## Usage

### Initialization

```go
// Initialize dependencies
repo := repositories.NewCourseRepository()
aiAnalyzer := ai.NewGeminiAnalyzer()

// Create and start the statistics service
statisticsService := queue.NewStatisticsService(repo, aiAnalyzer)
statisticsService.Start()

// Initialize course handler with the service
courseHandler := course.NewCourseHandler(
    repo,
    notification,
    aiAnalyzer,
    metricsClient,
    statisticsService,
)

// Don't forget to stop the service on shutdown
defer statisticsService.Stop()
```

### Enqueueing Tasks

The course handler automatically enqueues statistics calculation tasks when:
- A submission is created or updated
- A submission is deleted
- A submission is graded

Tasks are processed asynchronously by background workers.

## Configuration

Default configuration:
- **Workers**: 3 worker goroutines
- **Buffer Size**: 100 tasks
- **Max Retries**: 3 retries per task
- **Retry Delay**: Exponential backoff (1s, 2s, 3s)

## Benefits

1. **Non-blocking**: HTTP requests return immediately without waiting for statistics calculations
2. **Resilient**: Automatic retry mechanism for failed tasks
3. **Scalable**: Multiple workers can process tasks in parallel
4. **Graceful Shutdown**: Proper cleanup when the application stops
5. **Monitoring**: Queue size tracking and logging for observability

## Error Handling

- Tasks that fail are automatically retried up to the configured maximum
- Failed tasks after max retries are logged and discarded
- Worker failures don't affect other workers or the main application
- Graceful degradation: if the queue is full, errors are returned to prevent memory issues

## Monitoring

The system provides:
- Task enqueueing/processing logs
- Queue size metrics via `GetQueueSize()`
- Worker lifecycle logging
- Error and retry logging
