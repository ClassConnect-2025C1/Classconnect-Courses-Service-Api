package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"templateGo/controller"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	// Setup router
	mux := controller.SetupRoutes()

	// Create a test server (no need to specify port)
	server := httptest.NewServer(mux)
	defer server.Close() // This ensures the server is shut down after the test

	// Test "/healthcheck" endpoint using the server's URL
	resp, err := http.Get(server.URL + "/healthcheck")
	if err != nil {
		t.Fatalf("Error making GET request to '/healthcheck': %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for '/healthcheck', got: %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got: %s", contentType)
	}

	// Check response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response from '/healthcheck': %v", err)
	}
	expected := `{"status":"ok"}`
	if string(bodyBytes) != expected {
		t.Errorf("Response doesn't match, expected: %s, got: %s", expected, string(bodyBytes))
	}
}
