package reqres

import (
	"tanya_dokter_app/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GlobalCategorySpecialistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Status      int    `json:"status"`
}

func (request GlobalCategorySpecialistRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Code, validation.Required),
	)
}

type GlobalCategorySpecialistResponse struct {
	models.CustomGormModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Status      int    `json:"status"`
}

type GlobalCategorySpecialistBasePath struct {
	ID       int    `json:"id"`
	BasePath string `json:"base_path"`
}

type GlobalCategorySpecialistUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Code        string `json:"code"`
}
