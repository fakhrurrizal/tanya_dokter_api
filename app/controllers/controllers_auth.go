package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/repository"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

func HashID(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SignUp godoc
// @Summary SignUp
// @Description SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body reqres.SignUpRequest true "SignUp user"
// @Success 200
// @Router /v1/auth/signup [post]
func SignUp(c echo.Context) error {

	var request reqres.SignUpRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	utils.StripTagsFromStruct(&request)

	if err := request.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	_, err := repository.GetUserByEmail(strings.ToLower(request.Email))
	if err == nil {
		return c.JSON(http.StatusBadRequest, utils.Respond(http.StatusBadRequest, "bad request", "email sudah terdaftar"))
	}

	inputUser := reqres.GlobalUserRequest{
		Fullname: request.Fullname,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		Address:  request.Address,
		Gender:   request.Gender,
		Avatar:   request.Avatar,
		ZipCode:  request.ZipCode,
		Village:  request.Village,
		District: request.District,
		City:     request.City,
		Province: request.Province,
		Country:  request.Country,
	}

	if request.RoleID == 0 {
		inputUser.RoleID = 2
	}

	_, err = repository.CreateUser(1, true, &inputUser, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewBadRequestError([]map[string]interface{}{
			{
				"field": "Email",
				"error": err.Error(),
			},
		}))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Registration Successful",
	})

}

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Tags Auth
// @Accept json
// @Produce json
// @Param signin body reqres.SignInRequest true "SignIn user"
// @Success 200
// @Router /v1/auth/signin [post]
func SignIn(c echo.Context) error {

	var req reqres.SignInRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	user, accessToken, err := repository.SignIn(req.Email, req.Password)
	if err != nil {
		return c.JSON(400, map[string]interface{}{
			"status": 400,
			"error":  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user":         user,
			"access_token": accessToken,
			"expiration":   time.Now().Add(time.Hour * 72).Format("2006-01-02 15:04:05"),
		},
	})
}

// ResendEmailVerification godoc
// @Summary ResendEmail User
// @Description ResendEmail user
// @Tags Auth
// @Accept json
// @Produce json
// @Param Update body reqres.EmailRequest true "valid email"
// @Success 200
// @Router /v1/auth/resend-email-verification [post]
// @Security ApiKeyAuth
func ResendEmailVerification(c echo.Context) error {

	var req reqres.EmailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&req)

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	user, _ := repository.GetUserByEmail(strings.ToLower(req.Email))
	if user.EmailVerifiedAt.Time.IsZero() {
		data, err := repository.GetVerificationToken(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, utils.NewNotFoundError(err))
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": 200,
			"data":   data,
		})
	}

	return c.JSON(http.StatusForbidden, map[string]interface{}{
		"status":  403,
		"message": "Email has been verified",
	})
}

// EmailVerification godoc
// @Summary Email Verification for User
// @Description Email Verification for User
// @Tags Auth
// @Accept json
// @Produce json
// @Param Update body reqres.TokenRequest true "fill with valid token"
// @Success 200
// @Router /v1/auth/email-verification [post]
// @Security ApiKeyAuth
func EmailVerification(c echo.Context) error {
	var req reqres.TokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&req)

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	data, err := repository.EmailVerification(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewUnprocessableEntityError(err.Error()))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 200,
		"data":   data,
	})
}

// GetSignInUser godoc
// @Summary Get Sign In User
// @Description Get Sign In User
// @Tags Auth
// @Produce json
// @Success 200
// @Router /v1/auth/user [get]
// @Security ApiKeyAuth
// @Security JwtToken
func GetSignInUser(c echo.Context) error {

	id := c.Get("user_id").(int)
	user, err := repository.GetUserByIDPlain(id)
	if err != nil {
		return c.JSON(500, utils.Respond(500, err, "Failed to get user data"))
	}

	var roles models.GlobalRole

	var data reqres.GlobalUserAuthResponse
	data.ID = int(user.ID)

	data.Fullname = user.Fullname
	data.Avatar = user.Avatar
	data.Email = user.Email
	data.Phone = user.Phone
	data.Address = user.Address
	data.Village = user.Village
	data.District = user.District
	data.City = user.City
	data.Province = user.Province
	data.Country = user.Country
	data.ZipCode = user.ZipCode
	data.Status = user.Status
	data.Gender = user.Gender

	if user.EmailVerifiedAt.Time.IsZero() {
		data.EmailVerified = false
	} else {
		data.EmailVerified = true
	}

	if user.RoleID > 0 {
		roles, _ = repository.GetRoleByIDPlain(user.RoleID)
		data.Role = reqres.GlobalIDNameResponse{
			ID:   int(roles.ID),
			Name: roles.Name,
		}
	}

	return c.JSON(200, map[string]interface{}{
		"status":  200,
		"data":    data,
		"message": "Success to get user data",
	})
}

// ForgotPassword godoc
// @Summary Send Token Reset Password
// @Description Send Token Reset Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body reqres.EmailRequest true "Send token to email for reset password"
// @Success 200
// @Router /v1/auth/forgot-password [post]
// @Security ApiKeyAuth
func ForgotPassword(c echo.Context) error {
	var req reqres.EmailRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&req)

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	rand.Seed(time.Now().UnixNano())
	pin := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := repository.SaveResetRequest(req.Email, pin, expiresAt); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to save reset request",
		})
	}

	if err := repository.SendResetPassword(pin, req.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to send reset password email",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Reset password PIN has been sent",
	})
}

// ResetPassword godoc
// @Summary Reset User Password
// @Description Reset User Password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Update body reqres.ResetPasswordRequest true "body to update password"
// @Success 200
// @Router /v1/auth/reset-password [post]
// @Security ApiKeyAuth
func ResetPassword(c echo.Context) error {
	var req reqres.ResetPasswordRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}
	utils.StripTagsFromStruct(&req)

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := repository.ValidatePin(req.Email, req.Pin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// Reset password melalui repository
	if err := repository.UpdatePassword(req.Email, req.NewPassword); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Failed to reset password",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "Password has been reset",
	})
}
