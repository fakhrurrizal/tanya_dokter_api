package models

type GlobalRole struct {
	CustomGormModel
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Status      int    `json:"status" gorm:"column:status"`
}
