# AI Integration with Google Gemini

## 🤖 Overview

Este módulo integra **Google Gemini 2.0 Flash** como servicio de inteligencia artificial para proporcionar análisis automatizado y asistencia en el proceso educativo de la plataforma ClassConnect.

## 🔧 Architecture

### Files Structure
```
internal/handlers/ai/
├── gemini.go                   # Interface definition (FeedbackAnalyzer)
├── gemini_implementation.go    # Gemini API implementation
└── README.md                   # This documentation
```

### Interface Design
El sistema utiliza una interfaz `FeedbackAnalyzer` que permite cambiar fácilmente entre diferentes proveedores de IA sin afectar el resto del código:

```go
type FeedbackAnalyzer interface {
    GenerateCourseFeedbackAnalysis(courseTitle string, feedbacks []model.CourseFeedback) (string, error)
    GenerateGradeAndFeedback(assignmentDescription string, submissionFiles []model.SubmissionFile) (int, string, error)
    GenerateUserFeedbackAnalysis(feedbacks []model.UserFeedback) (string, error)
    GenerateCourseSuggestionsBasedOnStats(lastGradeTendency string, lastSubmissionRateTendency string, averageGrade float64) (string, error)
}
```

## 🚀 AI Provider: Google Gemini

### Model Used
- **Model**: `gemini-2.0-flash`
- **API**: Google Generative AI REST API
- **Multimodal**: Supports text and PDF document analysis

### Authentication
```go
// API Key configuration
apiKey := os.Getenv("GOOGLE_GEMINI_API_KEY")
```

## 🎯 AI Use Cases

### 1. **Automatic Grading and Feedback** 📝
**Endpoint**: `GET /{course_id}/assignment/{assignment_id}/submission/{submission_id}/ai-grade`

**What it does**:
- Analyzes student submissions (PDF files)
- Generates numerical grade (0-100)
- Provides detailed feedback and improvement suggestions

**Input**:
- Assignment description
- Student submission files (PDF format only)

**Output**:
```json
{
  "data": {
    "grade": 85,
    "feedback": "Excellent work on the theoretical concepts. The implementation shows good understanding of the algorithms. Consider adding more comments to improve code readability and include edge case handling."
  }
}
```

**Implementation**:
- Uses multimodal capabilities to analyze PDF content
- Processes both text and visual elements in documents
- Provides contextual feedback based on assignment requirements

### 2. **Course Feedback Analysis** 📊
**Endpoint**: `GET /{course_id}/ai-feedback-analysis`

**What it does**:
- Analyzes all feedback received for a course
- Identifies common themes and patterns
- Calculates average ratings and sentiment
- Provides insights for course improvement

**Input**:
- Course title
- Collection of student feedback (ratings + comments)

**Output**:
- Comprehensive analysis of student sentiment
- Key strengths identified by students
- Areas for improvement
- Average rating summary

**Sample Analysis**:
```text
The average rating for 'Introduction to Programming' is 4.2. Students consistently praise the clear explanations and practical examples. The main strengths include well-structured lessons and responsive instructor support. Areas for improvement include providing more challenging exercises for advanced students and extending deadline flexibility.
```

### 3. **User Performance Analysis** 👤
**Endpoint**: `GET /user/{user_id}/ai-feedback-analysis`

**What it does**:
- Analyzes feedback received by a specific student from multiple teachers
- Identifies student's strengths and weaknesses
- Provides personalized improvement recommendations

**Input**:
- Collection of feedback from different courses/teachers
- Student performance data

**Output**:
- Personalized analysis of student performance
- Cross-course pattern recognition
- Targeted improvement suggestions

### 4. **Course Statistics-Based Suggestions** 📈
**Function**: `GenerateCourseSuggestionsBasedOnStats`

**What it does**:
- Analyzes course performance trends
- Provides data-driven recommendations for course improvements
- Suggests actions based on grade and submission rate tendencies

**Input**:
- Grade tendency ("crescent", "decrescent", "stable")
- Submission rate tendency
- Average grade

**Output**:
- Strategic suggestions for course optimization
- Actionable recommendations based on data patterns

## ⚙️ Technical Implementation

### API Integration
```go
// Client creation
client, err := genai.NewClient(ctx, &genai.ClientConfig{
    APIKey:  g.APIKey,
    Backend: genai.BackendGeminiAPI,
})

// Content generation
result, err := client.Models.GenerateContent(
    ctx,
    "gemini-2.0-flash",
    contents,
    nil,
)
```

### Error Handling
- Comprehensive error handling for API failures
- Graceful degradation when AI service is unavailable
- Input validation for file types and content

### Performance Considerations
- Asynchronous processing for non-critical AI tasks
- File type restrictions (PDF only for grading)
- Response caching could be implemented for repeated analyses

## 🔗 Integration Points

### 1. Course Handlers Integration
The AI analyzer is injected into course handlers:

```go
courseHandler := course.NewCourseHandler(courseRepo, notificationClient, aiAnalyzer, ddMetrics, statisticsService)
```

### 2. Statistics Service Integration
AI suggestions are integrated with the queue-based statistics system:
- Course statistics calculations trigger AI-based suggestions
- Results are stored alongside statistical data
- Provides context-aware recommendations

### 3. Submission Processing
AI grading is available as an optional feature:
- Teachers can get AI-generated grades as suggestions
- Maintains human oversight in the grading process
- Supports teacher decision-making with AI insights

## 📋 API Response Formats

### Grading Response
```json
{
  "data": {
    "grade": 85,
    "feedback": "Detailed feedback text..."
  }
}
```

### Analysis Response
```json
{
  "data": {
    "analysis": "Comprehensive analysis text..."
  }
}
```

## 🛡️ Security & Privacy

### Data Handling
- Student submissions are processed securely
- No persistent storage of sensitive data in AI service
- API communications use HTTPS

### Privacy Considerations
- Student data is processed only for educational analysis
- No personal information is stored by the AI provider
- Results are anonymous and aggregated when possible

## 📊 Monitoring & Metrics

### Key Metrics to Track
- AI API response times
- Success/failure rates
- Grade accuracy (compared to human grading)
- User satisfaction with AI feedback
- API usage costs

### Error Monitoring
- API failures and timeouts
- Invalid file format attempts
- Malformed responses from AI service

## 🔧 Configuration

### Environment Variables (Recommended for Production)
```bash
GEMINI_API_KEY=your_api_key_here
GEMINI_MODEL=gemini-2.0-flash
AI_ENABLE=true
AI_TIMEOUT=30s
```

### Feature Flags
The AI features can be enabled/disabled through configuration without code changes, allowing for gradual rollout and testing.

---

## 🎓 Educational Impact

The AI integration enhances the educational experience by:

- **Reducing Teacher Workload**: Automated initial grading and feedback
- **Consistent Feedback**: Standardized evaluation criteria across submissions
- **Immediate Insights**: Real-time analysis of course performance
- **Personalized Learning**: Individual student performance insights
- **Data-Driven Decisions**: Evidence-based course improvements

This AI integration represents a significant step towards intelligent, adaptive educational technology that supports both teachers and students in achieving better learning outcomes.
