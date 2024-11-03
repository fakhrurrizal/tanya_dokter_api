package reqres

import (
	"tanya_dokter_app/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GlobalRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Flag        string `json:"flag"`
	Status      int    `json:"status"`
}

func (request GlobalRoleRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}

type GlobalRoleResponse struct {
	models.CustomGormModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Flag        string `json:"flag"`
	Status      int    `json:"status"`
}

type GlobalRoleBasePath struct {
	ID       int    `json:"id"`
	BasePath string `json:"base_path"`
}

type GlobalRoleWithUsers struct {
	models.CustomGormModel
	Name        string              `json:"name"`
	Description string              `json:"description"`
	TotalUsers  int                 `json:"total_users"`
	Users       []GlobalUserForRole `json:"users"`
}

type GlobalRoleWithUsersResponse struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	TotalUsers  int                 `json:"total_users"`
	Users       []GlobalUserForRole `json:"users"`
}

type GlobalRoleUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}
