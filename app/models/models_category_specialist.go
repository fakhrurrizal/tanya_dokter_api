package models

import (
	"time"

	"github.com/guregu/null"
)

type GlobalCategorySpecialist struct {
	CustomGormModel
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Image       string    `json:"image" gorm:"column:image"`
	Code        string    `json:"code" gorm:"column:code"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   null.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt   null.Time `json:"-" gorm:"column:deleted_at"`
	Status      int       `json:"status" gorm:"column:status"`
}
