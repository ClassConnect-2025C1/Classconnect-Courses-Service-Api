package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"templateGo/internal/model"
)

// GeminiAnalyzer implements FeedbackAnalyzer using Google's Gemini API
type GeminiAnalyzer struct {
	APIKey string
}

// NewGeminiAnalyzer creates a new instance of GeminiAnalyzer
func NewGeminiAnalyzer() *GeminiAnalyzer {
	// Get API key from environment
	// apiKey := os.Getenv("GEMINI_API_KEY")
	apiKey := "AIzaSyCGf5mrU_9zlsOg538SsjJSeq1yIyyLXDc"

	return &GeminiAnalyzer{
		APIKey: apiKey,
	}
}

// AnalyzeFeedback analyzes course feedback using the Gemini API
func (g *GeminiAnalyzer) AnalyzeFeedback(courseTitle string, feedbacks []model.CourseFeedback) (string, error) {
	// Format the feedback for the Gemini API
	feedbackText := formatFeedbackForAnalysis(courseTitle, feedbacks)

	// Call Gemini API
	return g.callGeminiAPI(feedbackText)
}

// formatFeedbackForAnalysis formats the feedback data into text format for Gemini
func formatFeedbackForAnalysis(courseTitle string, feedbacks []model.CourseFeedback) string {
	// Calculate the average rating
	// Calculate the average rating
	var totalRating int

	for _, feedback := range feedbacks {
		totalRating += feedback.Rating
	}

	averageRating := fmt.Sprintf("%.2f", float64(totalRating)/float64(len(feedbacks)))

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		"You are analyzing course feedbacks for '%s'. Your task is to provide a short and clear summary of the most common themes mentioned by students. First tell the average rating which is '%s', then identify key strengths and areas where the course can improve, considering the ratings are on a scale from 1 to 5. Output strictly plain text. Do not use lists, bullet points, bold text, markdown, or any kind of formatting. Keep it as a simple paragraph.",
		courseTitle,
		averageRating,
	))

	for i, feedback := range feedbacks {
		builder.WriteString(fmt.Sprintf("Feedback %d:\n", i+1))
		builder.WriteString(fmt.Sprintf("Rating: %d/100\n", feedback.Rating))
		if feedback.Summary != "" {
			builder.WriteString(fmt.Sprintf("Summary: %s\n", feedback.Summary))
		}
		if feedback.Comment != "" {
			builder.WriteString(fmt.Sprintf("Comment: %s\n", feedback.Comment))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

// callGeminiAPI calls the Google Gemini API to get an analysis of the feedback
func (g *GeminiAnalyzer) callGeminiAPI(feedbackText string) (string, error) {
	if g.APIKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY not available")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", g.APIKey)

	requestBody := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"text": feedbackText,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error calling Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("error unmarshaling response body: %w", err)
	}

	// Extract the generated text from the response
	candidates, ok := responseData["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	candidate := candidates[0].(map[string]any)
	content, ok := candidate["content"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	parts, ok := content["parts"].([]any)
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	part := parts[0].(map[string]any)
	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format from Gemini API")
	}

	return text, nil
}
