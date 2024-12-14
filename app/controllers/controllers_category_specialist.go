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

// GetCategorySpecialists godoc
// @Summary Create CategorySpecialists
// @Description Create New CategorySpecialists
// @Tags CategorySpecialist
// @Produce json
// @Param Body body reqres.GlobalCategorySpecialistRequest true "Create body"
// @Success 200
// @Router /v1/category-specialist [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateCategorySpecialist(c echo.Context) error {
	var input reqres.GlobalCategorySpecialistRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := input.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(400, utils.NewInvalidInputError(errVal))
	}

	CategorySpecialist, err := repository.CreateCategorySpecialist(&input)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to create CategorySpecialist"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    CategorySpecialist,
		"message": "Success to create CategorySpecialist",
	})
}

// GetCategorySpecialists godoc
// @Summary Get All CategorySpecialists With Pagination
// @Description Get All CategorySpecialists With Pagination
// @Tags CategorySpecialist
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param sort query string false "sort (ASC/DESC)"
// @Param order query string false "order by (default: id)"
// @Param status query boolean false "status (true (active) or false (inactive))"
// @Param created_at_margin_top query string false "created_at_margin_top (format: 2006-01-02)"
// @Param created_at_margin_bottom query string false "created_at_margin_top (format: 2006-01-02)"
// @Param code query string false "code (string)"
// @Produce json
// @Success 200
// @Router /v1/category-specialist [get]
// @Security ApiKeyAuth
func GetCategorySpecialists(c echo.Context) error {

	param := utils.PopulatePaging(c, "status")
	data := repository.GetCategorySpecialists(param)

	return c.JSON(http.StatusOK, data)
}

// GetCategorySpecialistByID godoc
// @Summary Get Single CategorySpecialist
// @Description Get Single CategorySpecialist
// @Tags CategorySpecialist
// @Param id path integer true "ID"
// @Produce json
// @Success 200
// @Router /v1/category-specialist/{id} [get]
// @Security ApiKeyAuth
func GetCategorySpecialistByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetCategorySpecialistByID(id)
	if err != nil {
		return c.JSON(404, utils.Respond(404, err, "Record not found"))
	}
	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get",
	})
}

// UpdateCategorySpecialistByID godoc
// @Summary Update Single CategorySpecialist by ID
// @Description Update Single CategorySpecialist by ID
// @Tags CategorySpecialist
// @Produce json
// @Param id path integer true "ID"
// @Param Body body reqres.GlobalCategorySpecialistUpdateRequest true "Update body"
// @Success 200
// @Router /v1/category-specialist/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateCategorySpecialistByID(c echo.Context) error {
	var input reqres.GlobalCategorySpecialistUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetCategorySpecialistByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Record not found"))
	}

	if input.Name != "" {
		data.Name = input.Name
	}
	if input.Description != "" {
		data.Description = input.Description
	}

	if input.Code != "" {
		data.Code = input.Code
	}

	data.Status = input.Status

	dataUpdate, err := repository.UpdateCategorySpecialist(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to update"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    dataUpdate,
		"message": "Success to update",
	})
}

// DeleteCategorySpecialistByID godoc
// @Summary Delete Single CategorySpecialist by ID
// @Description Delete Single CategorySpecialist by ID
// @Tags CategorySpecialist
// @Produce json
// @Param id path integer true "ID"
// @Success 200
// @Router /v1/category-specialist/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteCategorySpecialistByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetCategorySpecialistByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Record not found"))
	}

	_, err = repository.DeleteCategorySpecialist(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to delete"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to delete",
	})
}
