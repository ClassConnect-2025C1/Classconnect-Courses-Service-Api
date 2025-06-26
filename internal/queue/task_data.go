package queue

// CourseStatisticsTaskData represents data for course statistics calculation task
type CourseStatisticsTaskData struct {
	CourseID  uint   `json:"course_id"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}

// UserCourseStatisticsTaskData represents data for user course statistics calculation task
type UserCourseStatisticsTaskData struct {
	CourseID  uint   `json:"course_id"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
}
