package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// DatadogLogger handles sending logs directly to Datadog HTTP API
type DatadogLogger struct {
	APIKey     string
	Source     string
	Service    string
	HostName   string
	Site       string // Add this field
	HTTPClient *http.Client
}

// LogEntry represents a single log entry to be sent to Datadog
type LogEntry struct {
	Message    string                 `json:"message"`
	Status     string                 `json:"status,omitempty"` // error, warning, info, etc.
	Service    string                 `json:"service"`
	Hostname   string                 `json:"hostname"`
	Source     string                 `json:"ddsource"`
	Tags       []string               `json:"ddtags,omitempty"`
	Timestamp  int64                  `json:"timestamp,omitempty"` // Unix timestamp in milliseconds
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// NewDatadogLogger creates a new Datadog logger
func NewDatadogLogger(apiKey string) *DatadogLogger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Get site from environment or use provided default
	site := os.Getenv("DATADOG_SITE")
	if site == "" {
		site = "us5.datadoghq.com" // Default to US5 site
	}

	return &DatadogLogger{
		APIKey:     apiKey,
		Source:     "go",
		Service:    "classconnect-courses-api",
		HostName:   hostname,
		Site:       site,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Info logs an informational message
func (d *DatadogLogger) Info(message string, attributes map[string]interface{}, tags []string) error {
	return d.Log(message, "info", attributes, tags)
}

// Error logs an error message
func (d *DatadogLogger) Error(message string, attributes map[string]interface{}, tags []string) error {
	return d.Log(message, "error", attributes, tags)
}

// Warn logs a warning message
func (d *DatadogLogger) Warn(message string, attributes map[string]interface{}, tags []string) error {
	return d.Log(message, "warning", attributes, tags)
}

// Log sends a log message to Datadog
func (d *DatadogLogger) Log(message, status string, attributes map[string]interface{}, tags []string) error {
	entry := LogEntry{
		Message:    message,
		Status:     status,
		Service:    d.Service,
		Hostname:   d.HostName,
		Source:     d.Source,
		Tags:       tags,
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Attributes: attributes,
	}

	return d.SendLogs([]LogEntry{entry})
}

// SendLogs sends multiple log entries to Datadog in a single request
func (d *DatadogLogger) SendLogs(logs []LogEntry) error {
	payload, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("error marshaling logs: %w", err)
	}

	fmt.Printf("Sending logs to Datadog: %s\n", string(payload))

	// Use the configured site in the URL
	url := fmt.Sprintf("https://http-intake.logs.%s/api/v2/logs", d.Site)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.APIKey)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending logs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error from Datadog API: status code %d, body: %s",
			resp.StatusCode, string(body))
	}

	return nil
}
