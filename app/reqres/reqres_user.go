package reqres

import (
	"tanya_dokter_app/app/models"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type GlobalUserRequest struct {
	Avatar       string `json:"avatar"`
	Fullname     string `json:"fullname" validate:"required"`
	Gender       string `json:"gender"`
	Email        string `json:"email" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Village      string `json:"village"`
	District     string `json:"district"`
	City         string `json:"city"`
	Province     string `json:"province"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
	RoleID       int    `json:"role_id"`
	Status       int    `json:"status"`
	AutoVerified bool   `json:"auto_verified"`
}

func (request GlobalUserRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
		validation.Field(&request.Fullname, validation.Required),
	)
}

type GlobalUserUpdateRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Fullname     string `json:"fullname"`
	Phone        string `json:"phone"`
	Status       int    `json:"status"`
	RoleID       int    `json:"role_id"`
	ReferralCode string `json:"referral_code"`
	Token        int    `json:"token"`
	UserParentID int    `json:"-"`
	Gender       string `json:"gender"`
}

func (request GlobalUserUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Fullname, validation.Required),
		validation.Field(&request.Email, validation.Required, is.Email),
	)
}

type GlobalUserProfileUpdateRequest struct {
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Village  string `json:"village"`
	District string `json:"district"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	ZipCode  string `json:"zip_code"`
	Gender   string `json:"gender"`
}

func (request GlobalUserProfileUpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Fullname, validation.Required),
		validation.Field(&request.Password, validation.Required),
		validation.Field(&request.Email, validation.Required),
	)
}

type GlobalUserResponse struct {
	models.CustomGormModel
	Avatar       string               `json:"avatar"`
	Fullname     string               `json:"fullname" validate:"required"`
	Email        string               `json:"email" validate:"required"`
	Password     string               `json:"password" validate:"required"`
	Phone        string               `json:"phone"`
	Address      string               `json:"address"`
	Village      string               `json:"village"`
	District     string               `json:"district"`
	City         string               `json:"city"`
	Province     string               `json:"province"`
	Country      string               `json:"country"`
	ZipCode      string               `json:"zip_code"`
	Role         GlobalIDNameResponse `json:"role"`
	AutoVerified bool                 `json:"auto_verified"`
	Gender       string               `json:"gender"`
	Status       int                  `json:"status"`
}

type GlobalUserForRole struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CustomGormModel struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	EncodedID string    `json:"encoded_id" gorm:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SignUpRequest struct {
	Avatar   string `json:"avatar"`
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Village  string `json:"village"`
	District string `json:"district"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
	ZipCode  string `json:"zip_code"`
	RoleID   int    `json:"role_id"`
	Gender   string `json:"gender"`
}

func (request SignUpRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Length(8, 30)),
		validation.Field(&request.Fullname, validation.Required, validation.Length(5, 50)),
		validation.Field(&request.Phone, validation.Length(7, 17).Error("Nomor telepon harus benar")),
		validation.Field(&request.Address, validation.Length(5, 100)),
		validation.Field(&request.Gender, validation.In("m", "f").Error("Gender harus m(Laki) atau f(perempuan)")),
	)
}

type NewUserEmailNotification struct {
	AppName       string
	AdminFullname string
	FrontendUrl   string
	Fullname      string
	Email         string
	Phone         string
	CompanyName   string
	ContactEmail  string
}
