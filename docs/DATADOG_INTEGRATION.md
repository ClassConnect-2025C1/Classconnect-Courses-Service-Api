# Datadog Integration

## 📊 Overview

Este proyecto integra **Datadog** como plataforma de observabilidad para monitoreo, logging y métricas en tiempo real. La implementación permite tracking completo del rendimiento de la aplicación, errores, y análisis de comportamiento del usuario.

## 🏗️ Architecture

### Files Structure
```
internal/
├── logger/
│   └── datadog_logger.go      # Logging to Datadog HTTP API
└── metrics/
    └── datadog_metrics.go     # Metrics to Datadog API
```

### Integration Points
```
service_controller.go → Creates Datadog clients → Injects into handlers → Sends logs/metrics
```

## 🔧 Implementation Details

### 1. **Datadog Logger** (`internal/logger/datadog_logger.go`)

#### Purpose
- Centralized logging system
- Direct HTTP API integration with Datadog
- Structured log entries with metadata

#### Configuration
```go
type DatadogLogger struct {
    APIKey     string           // Datadog API key
    Source     string           // Log source identifier ("go")
    Service    string           // Service name ("classconnect-courses-api")
    HostName   string           // Server hostname
    Site       string           // Datadog site (us5.datadoghq.com)
    HTTPClient *http.Client     // HTTP client for API calls
}
```

#### Log Entry Structure
```go
type LogEntry struct {
    Message    string                 `json:"message"`
    Status     string                 `json:"status"`        // error, warning, info
    Service    string                 `json:"service"`
    Hostname   string                 `json:"hostname"`
    Source     string                 `json:"ddsource"`
    Tags       []string               `json:"ddtags"`
    Timestamp  int64                  `json:"timestamp"`     // Unix timestamp (ms)
    Attributes map[string]interface{} `json:"attributes"`    // Custom metadata
}
```

#### Usage Examples
```go
// Info logging
ddLogger.Info("User successfully enrolled", map[string]any{
    "user_id": "12345",
    "course_id": 67,
    "enrollment_type": "standard"
}, []string{"enrollment", "success"})

// Error logging
ddLogger.Error("Database connection failed", map[string]any{
    "error": err.Error(),
    "retry_count": 3
}, []string{"database", "error"})
```

### 2. **Datadog Metrics** (`internal/metrics/datadog_metrics.go`)

#### Purpose
- Performance monitoring
- Business metrics tracking
- Real-time application health monitoring

#### Configuration
```go
type DatadogMetricsClient struct {
    APIKey     string        // Datadog API key
    Site       string        // Datadog site endpoint
    HTTPClient *http.Client  // HTTP client for API calls
}
```

#### Metric Types Supported
```go
// Counter metrics (incremental)
metricsClient.IncrementCounter("api.requests.total", []string{
    "endpoint:/course",
    "method:POST",
    "status:200"
})

// Custom metrics (gauge, count, etc.)
metricsClient.SendMetric("course.enrollment.count", 150.0, "gauge", []string{
    "course_id:123",
    "semester:fall2025"
})
```

## 🚀 Integration in Application

### Service Controller Setup
```go
func SetupRoutes(ddLogger *logger.DatadogLogger, ddMetrics *metrics.DatadogMetricsClient) *ServiceManager {
    // Datadog clients are injected into the service setup
    courseHandler := course.NewCourseHandler(courseRepo, notificationClient, aiAnalyzer, ddMetrics, statisticsService)
    
    // Middleware for automatic logging
    r.Use(func(c *gin.Context) {
        // Process request
        c.Next()
        
        // Log request details
        if ddLogger != nil {
            status := c.Writer.Status()
            path := c.Request.URL.Path
            method := c.Request.Method
            
            attributes := map[string]any{
                "status":    status,
                "path":      path,
                "method":    method,
                "client_ip": c.ClientIP(),
            }
            
            if status >= 400 {
                ddLogger.Error(fmt.Sprintf("%s %s - %d", method, path, status), attributes, nil)
            } else {
                ddLogger.Info(fmt.Sprintf("%s %s - %d", method, path, status), attributes, nil)
            }
        }
    })
}
```

### Main Application Initialization
```go
func main() {
    // Initialize Datadog clients
    ddLogger := logger.NewDatadogLogger(datadogAPIKey)
    ddMetrics := metrics.NewDatadogMetricsClient(datadogAPIKey)
    
    // Send startup log
    ddLogger.Info("Application starting up", map[string]any{
        "version": "1.0.0",
        "environment": os.Getenv("ENV"),
    }, []string{"startup", "init"})
    
    // Setup routes with Datadog integration
    serviceManager := services.SetupRoutes(ddLogger, ddMetrics)
}
```

## 📈 What We Monitor

### 1. **HTTP Request Logging**
Every HTTP request is automatically logged with:
- **Request Method** (GET, POST, PUT, DELETE)
- **Request Path** (endpoint accessed)
- **Response Status** (200, 404, 500, etc.)
- **Client IP** (for security analysis)
- **Response Time** (performance tracking)

```json
{
  "message": "POST /course - 201",
  "status": "info",
  "service": "classconnect-courses-api",
  "hostname": "server-01",
  "ddsource": "go",
  "timestamp": 1751422950658,
  "attributes": {
    "status": 201,
    "path": "/course",
    "method": "POST",
    "client_ip": "192.168.1.100"
  }
}
```

### 2. **Application Events**
Business logic events are tracked:
- **User Enrollments**
- **Course Creation**
- **Assignment Submissions**
- **Grade Updates**
- **AI Analysis Requests**

### 3. **Error Tracking**
Comprehensive error monitoring:
- **Database Errors**
- **Authentication Failures**
- **API Integration Errors**
- **Validation Errors**
- **System Exceptions**

### 4. **Performance Metrics**
- **API Response Times**
- **Database Query Performance**
- **Queue Processing Times**
- **AI Service Response Times**
- **Resource Utilization**

## 🔧 Configuration

### Environment Variables
```bash
# Datadog Configuration
DATADOG_API_KEY=your_datadog_api_key_here
DATADOG_SITE=us5.datadoghq.com

# Application Configuration
SERVICE_NAME=classconnect-courses-api
ENVIRONMENT=production
VERSION=1.0.0
```

### Datadog Site Configuration
The implementation supports multiple Datadog sites:
- `us1.datadoghq.com` (US1)
- `us3.datadoghq.com` (US3)
- `us5.datadoghq.com` (US5) - **Default**
- `eu1.datadoghq.eu` (EU)
- `ap1.datadoghq.com` (AP)

## 📊 Datadog Dashboards

### Key Metrics to Monitor

#### 1. **API Performance Dashboard**
```
- api.requests.total (by endpoint, status)
- api.response_time.avg (by endpoint)
- api.errors.rate (by error type)
- api.throughput (requests per minute)
```

#### 2. **Business Metrics Dashboard**
```
- course.enrollments.total
- assignments.submissions.total
- ai.analyses.total
- user.activity.sessions
```

#### 3. **Infrastructure Dashboard**
```
- system.cpu.usage
- system.memory.usage
- database.connections.active
- queue.tasks.pending
```

#### 4. **Error Monitoring Dashboard**
```
- errors.by_endpoint
- errors.by_type
- database.errors.rate
- authentication.failures.rate
```

### Sample Queries
```sql
-- API Error Rate
sum:api.requests.total{status:error} / sum:api.requests.total{*}

-- Average Response Time by Endpoint
avg:api.response_time{*} by {endpoint}

-- Course Enrollment Trend
sum:course.enrollments.total{*}.as_rate()
```

## 🚨 Alerting Strategy

### Critical Alerts
1. **High Error Rate**: > 5% of requests return 5xx errors
2. **Slow Response Time**: API response time > 2 seconds
3. **Database Issues**: Connection failures or query timeouts
4. **Queue Backlog**: > 1000 pending tasks in queue

### Warning Alerts
1. **Increased Error Rate**: > 2% of requests return 4xx/5xx errors
2. **Performance Degradation**: Response time > 1 second
3. **High Memory Usage**: > 80% memory utilization
4. **AI Service Issues**: AI analysis failures > 10%

### Alert Channels
- **Slack**: Real-time notifications to development team
- **Email**: Critical alerts to on-call engineers
- **PagerDuty**: Escalation for production incidents

## 🔍 Log Analysis Examples

### Finding User Activity Patterns
```sql
-- User enrollment patterns
service:classconnect-courses-api "User successfully enrolled" 
| stats count by @attributes.course_id
```

### Error Investigation
```sql
-- Database connection errors
service:classconnect-courses-api status:error "Database connection failed"
| timeseries span:5m
```

### Performance Analysis
```sql
-- Slow endpoints
service:classconnect-courses-api @attributes.status:>=400
| stats avg(@duration) by @attributes.path
```

## 🔒 Security Considerations

### API Key Management
- **Environment Variables**: API keys stored securely in environment
- **Rotation**: Regular API key rotation policy
- **Access Control**: Limited permissions for API keys

### Data Privacy
- **PII Filtering**: Personal information excluded from logs
- **Data Retention**: Logs retained according to compliance requirements
- **Encryption**: All data transmitted to Datadog is encrypted (HTTPS)

### Network Security
- **Firewall Rules**: Outbound HTTPS access to Datadog endpoints
- **Rate Limiting**: Built-in rate limiting to prevent API abuse
- **Timeout Configuration**: Reasonable timeouts to prevent hanging requests

## 🚀 Benefits Achieved

### 1. **Operational Excellence**
- **Real-time Monitoring**: Immediate visibility into application health
- **Proactive Alerting**: Issues detected before users are affected
- **Root Cause Analysis**: Detailed logs help identify problem sources
- **Performance Optimization**: Data-driven performance improvements

### 2. **Business Intelligence**
- **User Behavior Analysis**: Understanding how users interact with the platform
- **Feature Usage Metrics**: Data on which features are most/least used
- **Performance Trends**: Long-term performance and usage trends
- **Capacity Planning**: Data-driven infrastructure scaling decisions

### 3. **Development Efficiency**
- **Faster Debugging**: Detailed error context speeds up problem resolution
- **Deployment Monitoring**: Real-time feedback on deployment health
- **Performance Regression Detection**: Automatic detection of performance issues
- **Code Quality Metrics**: Insights into code reliability and performance

## 📋 Maintenance Tasks

### Regular Tasks
1. **Log Retention Review**: Monitor log volume and costs
2. **Dashboard Updates**: Keep dashboards relevant to current features
3. **Alert Tuning**: Adjust alert thresholds based on baseline performance
4. **API Key Rotation**: Regular security maintenance

### Quarterly Reviews
1. **Cost Optimization**: Review Datadog usage and optimize costs
2. **Dashboard Cleanup**: Remove unused or outdated dashboards
3. **Metric Evaluation**: Assess if tracked metrics still provide value
4. **Alert Effectiveness**: Review alert accuracy and response times

---

## 🎯 Impact on ClassConnect

The Datadog integration provides ClassConnect with:

- **📊 Data-Driven Decisions**: Real metrics guide product development
- **🔍 Operational Visibility**: Complete transparency into system health
- **⚡ Faster Issue Resolution**: Detailed logs speed up debugging
- **📈 Performance Optimization**: Continuous monitoring drives improvements
- **🛡️ Proactive Problem Prevention**: Alerts catch issues before they impact users

This comprehensive observability foundation ensures ClassConnect can scale reliably while maintaining high performance and user satisfaction.
