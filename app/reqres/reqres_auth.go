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
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (request ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Token, validation.Required),
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
	ID            int                    `json:"id"`
	Email         string                 `json:"email"`
	Username      string                 `json:"username"`
	Fullname      string                 `json:"fullname"`
	Phone         string                 `json:"phone"`
	Status        int                    `json:"status"`
	Role          GlobalIDNameResponse   `json:"role,omitempty"`
	App           []GlobalIDNameResponse `json:"app,omitempty"`
	UserParent    GlobalIDNameResponse   `json:"user_parent,omitempty"`
	EmailVerified bool                   `json:"email_verified"`
}

type GlobalIDNameResponse struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Description string `json:"description,omitempty"`
	BaseIcon    string `json:"base_icon,omitempty"`
	Avatar      string `json:"image,omitempty"`
	// GeofencePointLimit int    `json:"geofence_point_limit,omitempty"`
}
