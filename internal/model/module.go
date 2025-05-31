package model

type Module struct {
	ID       uint   `gorm:"primaryKey" json:"id"` // Add this primary key
	CourseID uint   `gorm:"not null" json:"course_id"`
	Order    int    `gorm:"not null;default:0" json:"order"`
	Name     string `gorm:"not null" json:"name"`

	// Associations
	Course Course `gorm:"foreignKey:CourseID" json:"-"`
}

type Resource struct {
	ID       string `gorm:"primaryKey" json:"id"` // Add this primary key
	ModuleID uint   `gorm:"not null" json:"module_id"`
	Order    int    `gorm:"not null;default:0" json:"order"` // Add default value
	Type     string `gorm:"not null" json:"type"`
	URL      string `gorm:"not null" json:"url"`

	// Associations
	Module Module `gorm:"foreignKey:ModuleID;references:ID" json:"-"` // Fix relationship
}
