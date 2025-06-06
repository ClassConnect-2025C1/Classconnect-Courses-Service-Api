package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DatadogMetricsClient handles sending metrics to Datadog
type DatadogMetricsClient struct {
	APIKey     string
	Site       string
	HTTPClient *http.Client
}

// MetricPoint represents a single data point
type MetricPoint struct {
	Metric string          `json:"metric"`
	Points [][]interface{} `json:"points"` // [[timestamp, value], ...]
	Type   string          `json:"type"`   // gauge, count, etc.
	Tags   []string        `json:"tags,omitempty"`
	Host   string          `json:"host,omitempty"`
}

// MetricsPayload is the full payload sent to Datadog
type MetricsPayload struct {
	Series []MetricPoint `json:"series"`
}

// NewDatadogMetricsClient creates a new Datadog metrics client
func NewDatadogMetricsClient(apiKey string) *DatadogMetricsClient {
	// Get site from environment or use provided default
	site := os.Getenv("DATADOG_SITE")
	if site == "" {
		site = "us5.datadoghq.com" // Default to US5 site
	}

	return &DatadogMetricsClient{
		APIKey:     apiKey,
		Site:       site,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// IncrementCounter increments a counter metric
func (d *DatadogMetricsClient) IncrementCounter(metricName string, tags []string) error {
	return d.SendMetric(metricName, 1.0, "count", tags)
}

// SendMetric sends a single metric to Datadog
func (d *DatadogMetricsClient) SendMetric(metricName string, value float64, metricType string, tags []string) error {
	hostname, _ := os.Hostname()

	// Create the metric point
	point := MetricPoint{
		Metric: metricName,
		Points: [][]interface{}{{time.Now().Unix(), value}},
		Type:   metricType,
		Tags:   tags,
		Host:   hostname,
	}

	// Create the payload
	payload := MetricsPayload{
		Series: []MetricPoint{point},
	}

	return d.sendMetrics(payload)
}

// sendMetrics sends metrics payload to Datadog
func (d *DatadogMetricsClient) sendMetrics(payload MetricsPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling metrics: %w", err)
	}

	// Use the configured site in the URL
	url := fmt.Sprintf("https://api.%s/api/v1/series", d.Site)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.APIKey)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Datadog API: status code %d, body: %s",
			resp.StatusCode, string(body))
	}

	return nil
}
