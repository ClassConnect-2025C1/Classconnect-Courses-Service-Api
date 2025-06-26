package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"templateGo/internal/model"

	"google.golang.org/genai"
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

// GenerateGradeAndFeedback generates a grade and feedback for a submission using the Gemini API
func (g *GeminiAnalyzer) GenerateGradeAndFeedback(submissionDescription string, submissionFiles []model.SubmissionFile) (int, string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(
		ctx, &genai.ClientConfig{
			APIKey:  g.APIKey,
			Backend: genai.BackendGeminiAPI,
		})

	if err != nil {
		return 0, "", fmt.Errorf("error creating Gemini client: %w", err)
	}

	// Format the submission content and files for the Gemini API
	submissionText := "You are analyzing a student's submission for the following assignment:\n\n" + submissionDescription + "\n\nYour task is to provide a grade and feedback based on the content of the given submission files. The grade should be a number between 0 and 100, and the feedback should be a short paragraph explaining the grade and any suggestions for improvement.\n\n Answer strictly in plain text, in the following format:\n\n<grade>\n<feedback>\n\n"

	parts := []*genai.Part{}

	for _, file := range submissionFiles {
		if !strings.HasSuffix(file.Name, ".pdf") {
			return 0, "", fmt.Errorf("AI generated grade/feedback is only available for pdf files")
		}
		parts = append(parts,
			&genai.Part{
				InlineData: &genai.Blob{
					MIMEType: "application/pdf",
					Data:     file.Content,
				},
			},
		)
	}

	parts = append(parts, genai.NewPartFromText(submissionText))

	log.Println("Parts to be sent to Gemini API:", parts)

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		contents,
		nil,
	)
	if err != nil {
		return 0, "", fmt.Errorf("error generating content with Gemini API: %w", err)
	}

	// Extrat the grade and feedback from the result
	text := result.Text()

	grade_and_feedback := strings.SplitN(text, "\n", 2)
	if len(grade_and_feedback) < 2 {
		return 0, "", fmt.Errorf("unexpected response format from Gemini API: %s", text)
	}
	gradeText := strings.TrimSpace(grade_and_feedback[0])
	feedback := strings.TrimSpace(grade_and_feedback[1])
	grade, err := strconv.Atoi(gradeText)
	if err != nil {
		return 0, "", fmt.Errorf("error parsing grade from Gemini API response: %w", err)
	}
	if grade < 0 || grade > 100 {
		return 0, "", fmt.Errorf("grade out of range: %d", grade)
	}
	// Return the grade and feedback
	return grade, feedback, nil
}

// GenerateCourseFeedbackAnalysis analyzes course feedback using the Gemini API
func (g *GeminiAnalyzer) GenerateCourseFeedbackAnalysis(courseTitle string, feedbacks []model.CourseFeedback) (string, error) {
	// Format the feedback for the Gemini API
	feedbackText := formatCourseFeedbackForAnalysis(courseTitle, feedbacks)

	// Call Gemini API
	return g.callGeminiAPI(feedbackText)
}

// GenerateUserFeedbackAnalysis analyzes user feedback using the Gemini API
func (g *GeminiAnalyzer) GenerateUserFeedbackAnalysis(feedbacks []model.UserFeedback) (string, error) {
	// Format the feedback for the Gemini API
	feedbackText := formatUserFeedbackForAnalysis(feedbacks)

	// Call Gemini API
	return g.callGeminiAPI(feedbackText)
}

// formatCourseFeedbackForAnalysis formats the feedback data into text format for Gemini
func formatCourseFeedbackForAnalysis(courseTitle string, feedbacks []model.CourseFeedback) string {
	// Calculate the average rating
	var totalRating int

	for _, feedback := range feedbacks {
		totalRating += feedback.Rating
	}

	averageRating := fmt.Sprintf("%.2f", float64(totalRating)/float64(len(feedbacks)))

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		"You are analyzing course feedbacks for '%s'. Your task is to provide a short and clear summary of the most common themes mentioned by students. First tell the average rating which is '%s' (you don't have to recalculate it), then identify key strengths and areas where the course can improve considering that the ratings go from 1 to 5. Output strictly plain text. Do not use lists, bullet points, bold text, markdown, or any kind of formatting.",
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

// formatUserFeedbackForAnalysis formats the feedback data into text format for Gemini
func formatUserFeedbackForAnalysis(feedbacks []model.UserFeedback) string {
	// Calculate the average rating
	var totalRating uint

	for _, feedback := range feedbacks {
		totalRating += feedback.Rating
	}

	averageRating := fmt.Sprintf("%.2f", float64(totalRating)/float64(len(feedbacks)))

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(
		"You are analyzing feedbacks given to a student. Your task is to provide a short and clear summary of the comments and ratings from the teachers. First tell the average rating which is '%s' (you don't have to recalculate it), then identify key strengths and areas where the student can improve considering that the ratings go from 1 to 5. Output strictly plain text. Do not use lists, bullet points, bold text, markdown, or any kind of formatting.",
		averageRating,
	))

	for i, feedback := range feedbacks {
		builder.WriteString(fmt.Sprintf("Feedback %d:\n", i+1))
		builder.WriteString(fmt.Sprintf("Rating: %d/100\n", feedback.Rating))
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

func (g *GeminiAnalyzer) GenerateCourseSuggestionsBasedOnStats(lastGradeTendency string, lastSubmissionRateTendency string, averageGrade float64) (string, error) {
	if averageGrade == 0.0 {
		return "No stats to analyze", fmt.Errorf("assuming there is no submissions, average grade is 0")
	}
	// Format the input for the Gemini API
	inputText := fmt.Sprintf(
		"You are analyzing a course's statistics. The last grade tendency is '%s', the last submission rate tendency is '%s' and the last average grade is '%s'. Your task is to provide suggestions for improving the course based on these tendencies. Output strictly plain text. Do not use lists, bullet points, bold text, markdown, or any kind of formatting.",
		lastGradeTendency,
		lastSubmissionRateTendency,
		strconv.FormatFloat(averageGrade, 'f', 2, 64),
	)

	// Call Gemini API
	suggestions, err := g.callGeminiAPI(inputText)
	if err != nil {
		log.Printf("Error generating course suggestions: %v", err)
		return "Error generating suggestions", fmt.Errorf("error generating course suggestions: %w", err)
	}

	return suggestions, nil
}
