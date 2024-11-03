package models

import (
	"time"

	"github.com/guregu/null"
	"gorm.io/gorm"
)

type GlobalUser struct {
	CustomGormModel
	Email              string    `json:"email" gorm:"column:email"`
	Password           string    `json:"-" gorm:"column:password"`
	Fullname           string    `json:"fullname" gorm:"column:fullname"`
	Phone              string    `json:"phone" gorm:"column:phone"`
	Status             int       `json:"status" gorm:"column:status"`
	Gender             string    `json:"gender" gorm:"column:gender"`
	Avatar             string    `json:"avatar" gorm:"column:avatar"`
	Address            string    `json:"address" gorm:"column:address"`
	Village            string    `json:"village" gorm:"type: varchar(255)"`
	District           string    `json:"district" gorm:"type: varchar(255)"`
	City               string    `json:"city" gorm:"type: varchar(255)"`
	Province           string    `json:"province" gorm:"type: varchar(255)"`
	Country            string    `json:"country"`
	ZipCode            string    `json:"zip_code" gorm:"type: varchar(50)"`
	RoleID             int       `json:"role_id"`
	EmailVerifiedAt    null.Time `json:"-" gorm:"column:email_verified_at"`
	TwoFactorConfirmed bool      `json:"-" gorm:"column:two_factor_confirmed"`
	CreatedAt          time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt          null.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt          null.Time `json:"-" gorm:"column:deleted_at"`
}

type GlobalUserAllResponse struct {
	ID                 int                  `json:"id" gorm:"column:id;primaryKey:auto_increment"`
	Email              string               `json:"email" gorm:"column:email"`
	Fullname           string               `json:"fullname" gorm:"column:fullname"`
	Phone              string               `json:"phone" gorm:"column:phone"`
	Avatar             string               `json:"avatar" gorm:"column:avatar"`
	Status             int                  `json:"status" gorm:"column:status"`
	EmailVerifiedAt    null.Time            `json:"-" gorm:"column:email_verified_at"`
	TwoFactorConfirmed bool                 `json:"-" gorm:"column:two_factor_confirmed"`
	CreatedAt          time.Time            `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt          null.Time            `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt          null.Time            `json:"-" gorm:"column:deleted_at"`
	Role               GlobalIDNameResponse `json:"role"`
}

type GlobalIDNameResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GlobalUserDetail struct {
	ID     int    `json:"id" gorm:"column:id;primaryKey:auto_increment"`
	UserID int    `json:"user_id" gorm:"column:user_id"`
	Key    string `json:"key" gorm:"key"`
	Value  string `json:"value"`
}

type GlobalUserLogin struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey:auto_increment"`
	UserID      int       `json:"user_id" gorm:"column:user_id"`
	BearerToken string    `json:"bearer_token" gorm:"column:bearer_token"`
	IPAddress   string    `json:"ip_address" gorm:"column:ip_address"`
	UserAgent   string    `json:"user_agent" gorm:"column:user_agent"`
	CompanyID   int       `json:"company_id" gorm:"company_id"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt   null.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
	DeletedAt   null.Time `json:"-" gorm:"column:deleted_at"`
}

type CustomGormModel struct {
	ID        uint            `gorm:"primary_key" json:"id"`
	EncodedID string          `json:"encoded_id" gorm:"-"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
