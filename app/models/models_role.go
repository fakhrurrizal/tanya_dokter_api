package models

import (
	"time"

	"github.com/guregu/null"
)

type GlobalRole struct {
	CustomGormModel
	Name        string    `json:"name" gorm:"column:name"`
	Description string    `json:"description" gorm:"column:description"`
	Flag        string    `json:"flag" gorm:"column:flag"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   null.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt   null.Time `json:"-" gorm:"column:deleted_at"`
	Status      int       `json:"status" gorm:"column:status"`
}
