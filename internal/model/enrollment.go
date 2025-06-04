package model

// Enrollment representa la relación entre un usuario y un curso en la db
type Enrollment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   string `json:"user_id" gorm:"index"`
	CourseID uint   `json:"course_id" gorm:"index"`
	Favorite bool   `json:"favorite" gorm:"default:false"`
}
