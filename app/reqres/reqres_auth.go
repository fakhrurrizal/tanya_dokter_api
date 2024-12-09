package reqres

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request SignInRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required),
		validation.Field(&request.Password, validation.Required),
	)
}

type EmailRequest struct {
	Email string `json:"email"`
}

func (request EmailRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
	)
}

type ResetPasswordRequest struct {
	Email       string `gorm:"uniqueIndex"`
	Pin         string `gorm:"not null"`
	NewPassword string `json:"new_password"`
}

func (request ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Pin, validation.Required),
		validation.Field(&request.NewPassword, validation.Required),
	)
}

type TokenRequest struct {
	Email string `json:"email" validate:"required,email"`
	Pin   string `json:"pin" validate:"required,len=6,numeric"`
}

func (request TokenRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Pin, validation.Required),
	)
}

type GlobalUserAuthResponse struct {
	ID            int                  `json:"id"`
	Email         string               `json:"email"`
	Fullname      string               `json:"fullname"`
	Phone         string               `json:"phone"`
	Status        int                  `json:"status"`
	Gender        string               `json:"gender" `
	Avatar        string               `json:"avatar" `
	Address       string               `json:"address" `
	Village       string               `json:"village"`
	District      string               `json:"district"`
	City          string               `json:"city"`
	Province      string               `json:"province"`
	Country       string               `json:"country"`
	ZipCode       string               `json:"zip_code" gorm:"type: varchar(50)"`
	Role          GlobalIDNameResponse `json:"role,omitempty"`
	EmailVerified bool                 `json:"email_verified"`
}

type GlobalIDNameResponse struct {
	ID        int    `json:"id,omitempty"`
	EncodedID string `json:"encoded_id,omitempty"`
	Name      string `json:"name,omitempty"`
}
