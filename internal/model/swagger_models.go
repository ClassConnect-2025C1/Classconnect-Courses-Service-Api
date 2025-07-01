package model

// SuccessResponse represents a successful API response
// @Description Success response
type SuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
// @Description Error response
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Something went wrong"`
	Code    int    `json:"code" example:"400"`
}

// HealthCheckResponse represents the health check response
// @Description Health check response
type HealthCheckResponse struct {
	Status string `json:"status" example:"ok"`
}

// CourseResponse represents a course in API responses
// @Description Course information
type CourseResponse struct {
	ID          string `json:"id" example:"12345"`
	Title       string `json:"title" example:"Introduction to Programming"`
	Description string `json:"description" example:"Learn the basics of programming with Python"`
	CreatedBy   string `json:"createdBy" example:"teacher123"`
	CreatedAt   string `json:"created_at" example:"2023-01-15T10:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2023-01-15T10:00:00Z"`
}

// MembersList represents the members list for a course
// @Description Members list response
type MembersList struct {
	Data []Member `json:"data"`
}

// Member represents a course member
// @Description Course member information
type Member struct {
	Role  string `json:"role" example:"student"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
}

// AssignmentRequest represents the request body for creating an assignment
// @Description Request body for creating an assignment
type AssignmentRequest struct {
	Title       string `json:"title" binding:"required" example:"Programming Exercise 1"`
	Description string `json:"description" example:"Complete the following programming tasks"`
	DueDate     string `json:"dueDate" example:"2023-12-31T23:59:59Z"`
	TotalPoints int    `json:"totalPoints" example:"100"`
}

// SubmissionRequest represents the request body for submitting an assignment
// @Description Request body for submitting an assignment
type SubmissionRequest struct {
	Content     string   `json:"content" binding:"required" example:"Here is my solution to the assignment"`
	Attachments []string `json:"attachments,omitempty" example:"[\"file1.py\", \"file2.txt\"]"`
}

// GradeRequest represents the request body for grading a submission
// @Description Request body for grading a submission
type GradeRequest struct {
	Grade    float64 `json:"grade" binding:"required" example:"85.5"`
	Feedback string  `json:"feedback" example:"Good work! Consider improving the algorithm efficiency."`
}

// FeedbackRequest represents the request body for course feedback
// @Description Request body for course feedback
type FeedbackRequest struct {
	Content string  `json:"content" binding:"required" example:"This course was very helpful and well-structured"`
	Rating  float64 `json:"rating" example:"4.5"`
}

// Feedback represents feedback in API responses
// @Description Feedback information
type Feedback struct {
	ID        string  `json:"id" example:"feed123"`
	Content   string  `json:"content" example:"This course was very helpful"`
	Rating    float64 `json:"rating" example:"4.5"`
	CreatedBy string  `json:"createdBy" example:"student123"`
	CreatedAt string  `json:"createdAt" example:"2023-12-30T10:00:00Z"`
}

// FeedbackResponse represents feedback in API responses
// @Description Feedback information
type FeedbackResponse struct {
	ID        string  `json:"id" example:"feedback123"`
	Content   string  `json:"content" example:"Great course!"`
	Rating    float64 `json:"rating" example:"4.5"`
	CreatedBy string  `json:"createdBy" example:"student123"`
	CreatedAt string  `json:"createdAt" example:"2023-12-30T10:00:00Z"`
}

// ModuleRequest represents the request body for creating a module
// @Description Request body for creating a module
type ModuleRequest struct {
	Name string `json:"name" binding:"required" example:"Introduction to Variables"`
}

// ResourceRequest represents the request body for creating a resource
// @Description Request body for creating a resource
type ResourceRequest struct {
	Title   string `json:"title" binding:"required" example:"Python Variables Tutorial"`
	Type    string `json:"type" binding:"required" example:"video"`
	Content string `json:"content,omitempty" example:"Tutorial content here"`
	URL     string `json:"url,omitempty" example:"https://example.com/video"`
}

// AssignmentResponse represents an assignment in API responses
// @Description Assignment information
type AssignmentResponse struct {
	ID          string `json:"id" example:"assign123"`
	Title       string `json:"title" example:"Programming Assignment 1"`
	Description string `json:"description" example:"Implement a basic calculator"`
	DueDate     string `json:"dueDate" example:"2023-12-31T23:59:59Z"`
	TotalPoints int    `json:"totalPoints" example:"100"`
	CreatedAt   string `json:"createdAt" example:"2023-12-01T10:00:00Z"`
}
