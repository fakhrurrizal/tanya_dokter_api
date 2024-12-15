package models

type GlobalDataDrugs struct {
	CustomGormModel
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Image       string `json:"image" gorm:"column:image"`
	Code        string `json:"code" gorm:"column:code"`
	Usage        string `json:"usage" gorm:"column:usage"`
}
