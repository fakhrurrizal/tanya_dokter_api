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

// GetDataDrugs godoc
// @Summary Create DataDrugss
// @Description Create New DataDrugss
// @Tags DataDrugs
// @Produce json
// @Param Body body reqres.GlobalDataDrugsRequest true "Create body"
// @Success 200
// @Router /v1/data-drugs [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateDataDrugs(c echo.Context) error {
	var input reqres.GlobalDataDrugsRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := input.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(400, utils.NewInvalidInputError(errVal))
	}

	DataDrugs, err := repository.CreateDataDrugs(&input)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to create DataDrugs"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    DataDrugs,
		"message": "Success to create DataDrugs",
	})
}

// GetDataDrugs godoc
// @Summary Get All DataDrugss With Pagination
// @Description Get All DataDrugss With Pagination
// @Tags DataDrugs
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
// @Router /v1/data-drugs [get]
// @Security ApiKeyAuth
func GetDataDrugs(c echo.Context) error {

	param := utils.PopulatePaging(c, "status")
	data := repository.GetDataDrugs(param)

	return c.JSON(http.StatusOK, data)
}

// GetDataDrugsByID godoc
// @Summary Get Single DataDrugs
// @Description Get Single DataDrugs
// @Tags DataDrugs
// @Param id path integer true "ID"
// @Produce json
// @Success 200
// @Router /v1/data-drugs/{id} [get]
// @Security ApiKeyAuth
func GetDataDrugsByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetDataDrugsByID(id)
	if err != nil {
		return c.JSON(404, utils.Respond(404, err, "Record not found"))
	}
	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get",
	})
}

// UpdateDataDrugsByID godoc
// @Summary Update Single DataDrugs by ID
// @Description Update Single DataDrugs by ID
// @Tags DataDrugs
// @Produce json
// @Param id path integer true "ID"
// @Param Body body reqres.GlobalDataDrugsUpdateRequest true "Update body"
// @Success 200
// @Router /v1/data-drugs/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateDataDrugsByID(c echo.Context) error {
	var input reqres.GlobalDataDrugsUpdateRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetDataDrugsByIDPlain(id)
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

	if input.Image != "" {
		data.Image = input.Image
	}

	if input.Usage != "" {
		data.Usage = input.Usage
	}

	dataUpdate, err := repository.UpdateDataDrugs(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to update"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    dataUpdate,
		"message": "Success to update",
	})
}

// DeleteDataDrugsByID godoc
// @Summary Delete Single DataDrugs by ID
// @Description Delete Single DataDrugs by ID
// @Tags DataDrugs
// @Produce json
// @Param id path integer true "ID"
// @Success 200
// @Router /v1/data-drugs/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteDataDrugsByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := repository.GetDataDrugsByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(404, err, "Record not found"))
	}

	_, err = repository.DeleteDataDrugs(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to delete"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to delete",
	})
}
