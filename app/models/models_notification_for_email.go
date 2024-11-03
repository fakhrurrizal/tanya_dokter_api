package models

type NotificationForEmail struct {
	CustomGormModel
	UserID   int    `json:"user_id" gorm:"type: int8"`
	FullName string `json:"full_name" gorm:"type: varchar(255)"`
	Email    string `json:"email" gorm:"type: varchar(255)"`
	Status   bool   `json:"status" gorm:"type: bool"`
}
