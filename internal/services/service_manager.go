package services

import (
	"context"
	"net/http"
	"templateGo/internal/queue"
)

// ServiceManager manages the lifecycle of application services
type ServiceManager struct {
	statisticsService *queue.StatisticsService
	httpHandler       http.Handler
}

// NewServiceManager creates a new service manager
func NewServiceManager(statisticsService *queue.StatisticsService, httpHandler http.Handler) *ServiceManager {
	return &ServiceManager{
		statisticsService: statisticsService,
		httpHandler:       httpHandler,
	}
}

// Start starts all managed services
func (sm *ServiceManager) Start() {
	if sm.statisticsService != nil {
		sm.statisticsService.Start()
	}
}

// Stop stops all managed services gracefully
func (sm *ServiceManager) Stop() {
	if sm.statisticsService != nil {
		sm.statisticsService.Stop()
	}
}

// ServeHTTP implements http.Handler interface
func (sm *ServiceManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sm.httpHandler.ServeHTTP(w, r)
}

// Shutdown gracefully shuts down all services
func (sm *ServiceManager) Shutdown(ctx context.Context) error {
	// Stop statistics service
	sm.Stop()

	// If the HTTP handler supports graceful shutdown, call it here
	// For now, we just stop our services
	return nil
}
