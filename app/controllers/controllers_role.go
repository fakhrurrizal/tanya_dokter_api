package controllers

import (
	"net/http"
	"strconv"
	"tanya_dokter_app/app/repository"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

func CreateRole(c echo.Context) error {
	var input reqres.GlobalRoleRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := input.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(400, utils.NewInvalidInputError(errVal))
	}

	role, err := repository.CreateRole(&input)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to create role"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    role,
		"message": "Success to create role",
	})
}

// GetRoles godoc
// @Summary Get All Role With Pagination
// @Description Get All Role With Pagination
// @Tags Role
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param sort query string false "sort (ASC/DESC)"
// @Param order query string false "order by (default: id)"
// @Param status query boolean false "status (true (active) or false (inactive))"
// // @Param company_id query integer false "company_id (int)"
// @Param created_at_margin_top query string false "created_at_margin_top (format: 2006-01-02)"
// @Param created_at_margin_bottom query string false "created_at_margin_top (format: 2006-01-02)"
// @Param code query string false "code (string)"
// @Produce json
// @Success 200
// @Router /v1/role [get]
// @Security ApiKeyAuth
func GetRoles(c echo.Context) error {

	param := utils.PopulatePaging(c, "status")
	data := repository.GetRoles(param)

	return c.JSON(http.StatusOK, data)
}

// GetRoleByID godoc
// @Summary Get Single Role
// @Description Get Single Role
// @Tags Role
// @Param id path integer true "ID"
// @Produce json
// @Success 200
// @Router /v1/role/{id} [get]
// @Security ApiKeyAuth
func GetRoleByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetRoleByID(id)
	if err != nil {
		return c.JSON(404, utils.Respond(404, err, "Record not found"))
	}
	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get",
	})
}

// UpdateRoleByID godoc
// @Summary Update Single Role by ID
// @Description Update Single Role by ID
// @Tags Role
// @Produce json
// @Param id path integer true "ID"
// @Param Body body reqres.GlobalRoleUpdateRequest true "Update body"
// @Success 200
// @Router /v1/role/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateRoleByID(c echo.Context) error {
	var input reqres.GlobalRoleUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetRoleByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Record not found"))
	}

	if input.Name != "" {
		data.Name = input.Name
	}
	if input.Description != "" {
		data.Description = input.Description
	}

	data.Status = input.Status

	dataUpdate, err := repository.UpdateRole(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to update"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    dataUpdate,
		"message": "Success to update",
	})
}

// DeleteRoleByID godoc
// @Summary Delete Single Role by ID
// @Description Delete Single Role by ID
// @Tags Role
// @Produce json
// @Param id path integer true "ID"
// @Success 200
// @Router /v1/role/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteRoleByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetRoleByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Record not found"))
	}

	_, err = repository.DeleteRole(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to delete"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to delete",
	})
}
