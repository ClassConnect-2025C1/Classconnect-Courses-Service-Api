package model

// Enrollment representa la relación entre un usuario y un curso en la db
type Enrollment struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   string `json:"user_id" gorm:"index"`
	CourseID uint   `json:"course_id" gorm:"index"`
	Favorite bool   `json:"favorite" gorm:"default:false"`
}

// (1 - terminado) Modificar el response de los cursos aprobados para que devuelva el nombre del curso
// (2 - terminado) modificar el critterio de elegibilidad que sea una lista de strings (modificar el edit)

// (3) implementar los cursos favoritos (modificar el enrollment para que contenga un bool)
// (4) Chequear al momento del enrollment que el usuario cumpla con los criterios
// (4) chequear que el usuario inscripto realmente este inscripto al momento de ponerlo en favorito
