package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"tanya_dokter_app/app/middlewares"
	"tanya_dokter_app/app/repository"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

// CreateUser godoc
// @Summary Create User
// @Description Create New User
// @Tags User
// @Produce json
// @Param Body body reqres.GlobalUserRequest true "Create body"
// @Success 200
// @Router /v1/user [post]
// @Security ApiKeyAuth
// @Security JwtToken
func CreateUser(c echo.Context) error {

	var input reqres.GlobalUserRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	if err := input.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(400, utils.NewInvalidInputError(errVal))
	}

	userCheck, err := repository.GetUserByEmail(strings.ToLower(input.Email))
	if err == nil {
		return c.JSON(400, utils.Respond(400, err, "Email "+strings.ToLower(userCheck.Email)+" has been registered"))
	}

	data, err := repository.CreateUser(1, true, &input, 0)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to create"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Create Success",
	})
}

// GetUsers godoc
// @Summary Get All Users With Pagination
// @Description Get All Users With Pagination
// @Tags User
// @Param search query string false "search (string)"
// @Param page query integer false "page (int)"
// @Param limit query integer false "limit (int)"
// @Param sort query string false "sort (ASC/DESC)"
// @Param order query string false "order by (default: id)"
// @Param status query boolean false "status (true (active) or false (inactive))"
// @Param role_id query integer false "role_id (int)"
// @Param created_at_margin_top query string false "created_at_margin_top (format: 2006-01-02)"
// @Param created_at_margin_bottom query string false "created_at_margin_top (format: 2006-01-02)"
// @Produce json
// @Success 200
// @Router /v1/user [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetUsers(c echo.Context) error {

	roleID, _ := strconv.Atoi(c.QueryParam("role_id"))
	createdAtMarginTop := c.QueryParam("created_at_margin_top")
	createdAtMarginBottom := c.QueryParam("created_at_margin_bottom")

	param := utils.PopulatePaging(c, "status")
	data := repository.GetUsers(roleID, createdAtMarginTop, createdAtMarginBottom, param)

	return c.JSON(http.StatusOK, data)
}

// GetUserByID godoc
// @Summary Get Single User
// @Description Get Single User
// @Tags User
// @Param id path integer true "ID"
// @Produce json
// @Success 200
// @Router /v1/user/{id} [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetUserByID(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	data, err := repository.GetUserByID(id)
	if err != nil {
		return c.JSON(404, utils.Respond(404, err, "Record not found"))
	}
	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get",
	})
}

// DeleteUserByID godoc
// @Summary Delete Single User by ID
// @Description Delete Single User by ID
// @Tags User
// @Produce json
// @Param id path integer true "ID"
// @Param company_id query integer false "company_id (int)"
// @Success 200
// @Router /v1/user/{id} [delete]
// @Security ApiKeyAuth
// @Security JwtToken
func DeleteUserByID(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	if c.Get("user_id").(int) == id {
		return c.JSON(500, utils.Respond(500, "bad request", "You are not allowed to delete your own account"))
	}

	data, err := repository.GetUserByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to get"))
	}

	// response, err := repository.GetUserByID(int(data.ID), host.AppID, companyID, false, false)
	// if err != nil {
	// 	return c.JSON(404, utils.Respond(404, err, "Failed to get response"))
	// }

	_, err = repository.DeleteUser(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to delete"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to delete",
	})
}

// UpdateUserByID godoc
// @Summary Update Single User by ID
// @Description Update Single User by ID
// @Tags User
// @Produce json
// @Param id path integer true "ID"
// @Param Body body reqres.GlobalUserRequest true "Update body"
// @Param company_id query integer false "company_id (int)"
// @Success 200
// @Router /v1/user/{id} [put]
// @Security ApiKeyAuth
// @Security JwtToken
func UpdateUserByID(c echo.Context) error {

	var input reqres.GlobalUserRequest
	if err := c.Bind(&input); err != nil {
		return c.JSON(400, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&input)

	id, _ := strconv.Atoi(c.Param("id"))

	data, err := repository.GetUserByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to get"))
	}

	if input.Avatar != "" {
		data.Avatar = input.Avatar
	}
	if input.Fullname != "" {
		data.Fullname = input.Fullname
	}
	if input.Password != "" {
		data.Password = middlewares.BcryptPassword(input.Password)
	}
	if input.Email != "" {
		data.Email = strings.ToLower(input.Email)
	}
	if input.Phone != "" {
		data.Phone = input.Phone
	}
	if input.Address != "" {
		data.Address = input.Address
	}
	if input.Village != "" {
		data.Village = input.Village
	}
	if input.District != "" {
		data.District = input.District
	}

	if input.City != "" {
		data.City = input.City
	}
	if input.Province != "" {
		data.Province = input.Province
	}
	if input.Country != "" {
		data.Country = input.Country
	}
	if input.ZipCode != "" {
		data.ZipCode = input.ZipCode
	}

	if input.Status != 0 {
		data.Status = input.Status
	}

	if input.Gender != "" {
		data.Gender = input.Gender
	}

	// if input.RoleID != 0 {
	// 	data.RoleID = input.RoleID
	// }

	dataUpdate, err := repository.UpdateUser(data)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to update"))
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    dataUpdate,
		"message": "Success to update",
	})
}
