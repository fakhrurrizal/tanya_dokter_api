package reqres

import (
	"tanya_dokter_app/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GlobalDataDrugsRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Code        string `json:"code"`
	Usage       string `json:"usage"`
}

func (request GlobalDataDrugsRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Code, validation.Required),
	)
}

type GlobalDataDrugsResponse struct {
	models.CustomGormModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Code        string `json:"code"`
	Usage       string `json:"usage"`
}

type GlobalDataDrugsUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Code        string `json:"code"`
	Usage       string `json:"usage"`
}
