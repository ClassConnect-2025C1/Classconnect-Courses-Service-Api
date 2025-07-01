package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// TaskType represents the type of task to be executed
type TaskType string

const (
	TaskTypeCourseStatistics     TaskType = "course_statistics"
	TaskTypeUserCourseStatistics TaskType = "user_course_statistics"
	TaskTypeGlobalStatistics     TaskType = "global_statistics"
)

// Task represents a task to be executed
type Task struct {
	ID         string
	Type       TaskType
	Data       interface{}
	CreatedAt  time.Time
	Retries    int
	MaxRetries int
}

// TaskQueue represents a queue for processing tasks
type TaskQueue struct {
	tasks      chan Task
	workers    int
	workerPool sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	processor  TaskProcessor
	mu         sync.RWMutex
	running    bool
}

// TaskProcessor interface for processing tasks
type TaskProcessor interface {
	ProcessTask(task Task) error
}

// NewTaskQueue creates a new task queue
func NewTaskQueue(workers int, bufferSize int, processor TaskProcessor) *TaskQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &TaskQueue{
		tasks:     make(chan Task, bufferSize),
		workers:   workers,
		ctx:       ctx,
		cancel:    cancel,
		processor: processor,
		running:   false,
	}
}

// Start starts the task queue workers
func (tq *TaskQueue) Start() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if tq.running {
		return
	}

	tq.running = true

	for i := 0; i < tq.workers; i++ {
		tq.workerPool.Add(1)
		go tq.worker(i)
	}

	log.Printf("Task queue started with %d workers", tq.workers)
}

// Stop stops the task queue
func (tq *TaskQueue) Stop() {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if !tq.running {
		return
	}

	tq.running = false
	tq.cancel()
	close(tq.tasks)
	tq.workerPool.Wait()

	log.Println("Task queue stopped")
}

// EnqueueTask adds a task to the queue
func (tq *TaskQueue) EnqueueTask(task Task) error {
	tq.mu.RLock()
	defer tq.mu.RUnlock()

	if !tq.running {
		return fmt.Errorf("task queue is not running")
	}

	task.CreatedAt = time.Now()
	if task.MaxRetries == 0 {
		task.MaxRetries = 3 // Default max retries
	}

	select {
	case tq.tasks <- task:
		log.Printf("Task %s enqueued successfully", task.ID)
		return nil
	case <-tq.ctx.Done():
		return fmt.Errorf("task queue is shutting down")
	default:
		return fmt.Errorf("task queue is full")
	}
}

// worker processes tasks from the queue
func (tq *TaskQueue) worker(workerID int) {
	defer tq.workerPool.Done()

	log.Printf("Worker %d started", workerID)

	for {
		select {
		case task, ok := <-tq.tasks:
			if !ok {
				log.Printf("Worker %d stopped: channel closed", workerID)
				return
			}

			tq.processTask(workerID, task)

		case <-tq.ctx.Done():
			log.Printf("Worker %d stopped: context canceled", workerID)
			return
		}
	}
}

// processTask processes a single task with retry logic
func (tq *TaskQueue) processTask(workerID int, task Task) {
	log.Printf("Worker %d processing task %s (type: %s)", workerID, task.ID, task.Type)

	err := tq.processor.ProcessTask(task)
	if err != nil {
		log.Printf("Worker %d failed to process task %s: %v", workerID, task.ID, err)

		// Retry logic
		if task.Retries < task.MaxRetries {
			task.Retries++
			log.Printf("Retrying task %s (attempt %d/%d)", task.ID, task.Retries, task.MaxRetries)

			// Add a small delay before retry
			go func() {
				time.Sleep(time.Second * time.Duration(task.Retries))
				if err := tq.EnqueueTask(task); err != nil {
					log.Printf("Failed to re-enqueue task %s: %v", task.ID, err)
				}
			}()
		} else {
			log.Printf("Task %s failed after %d retries", task.ID, task.MaxRetries)
		}
	} else {
		log.Printf("Worker %d successfully processed task %s", workerID, task.ID)
	}
}

// GetQueueSize returns the current queue size
func (tq *TaskQueue) GetQueueSize() int {
	return len(tq.tasks)
}
