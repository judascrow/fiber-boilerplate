package models

// Role model
type Role struct {
	BaseModel
	Name        string `gorm:"unique;not null" json:"name" form:"name" query:"name"`
	Description string `gorm:"type:varchar(100);" json:"description"`
}
