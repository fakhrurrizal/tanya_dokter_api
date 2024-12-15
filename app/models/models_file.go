package models

type GlobalFile struct {
	CustomGormModel
	Token    string `json:"token" gorm:"type: varchar(255)"`
	UserID   int    `json:"user_id" gorm:"type: int8"`
	Filename string `json:"filename" gorm:"type: varchar(255)"`
	Path     string `json:"-" gorm:"type: varchar(255)"`
	FullUrl  string `json:"full_url" gorm:"-"`
}
