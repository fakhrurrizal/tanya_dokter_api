package repository

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"tanya_dokter_app/config"
	"time"

	"github.com/guregu/null"
)

func GetUserByEmail(email string) (response models.GlobalUser, err error) {
	err = config.DB.Where("email = ?", strings.ToLower(email)).First(&response).Error

	return
}

func CreateUser(status int, verification bool, data *reqres.GlobalUserRequest, userID int) (user models.GlobalUser, err error) {
	var userInput models.GlobalUser
	user, errUser := GetUserByEmail(data.Email)
	if errUser == nil {
		err = errors.New("email has been registered")
		return
	} else {
		passwordHashed := models.BcryptPassword(data.Password)
		var verifiedAt null.Time
		if data.Status != 0 {
			verifiedAt = null.TimeFrom(time.Now())
		}
		userInput = models.GlobalUser{
			Email:           data.Email,
			Password:        passwordHashed,
			Fullname:        data.Fullname,
			Phone:           data.Phone,
			Status:          status,
			RoleID:          data.RoleID,
			Address:         data.Address,
			Village:         data.Village,
			District:        data.District,
			City:            data.City,
			Province:        data.Province,
			Country:         data.Country,
			Avatar:          data.Avatar,
			EmailVerifiedAt: verifiedAt,
			ZipCode:         data.ZipCode,
			Gender:          data.Gender,
			CategoryID:      data.CategoryID,
			Experience:      data.Experience,
			Code:            data.Code,
		}
		err = config.DB.Create(&userInput).Error

	}

	user, _ = GetUserByEmail(data.Email)
	if verification {
		log.Println("Sending Email Condition")
		GetVerificationToken(&reqres.EmailRequest{Email: data.Email})
	}

	return user, err
}

func BuildUserResponse(data models.GlobalUser) (response reqres.GlobalUserResponse) {

	var roles models.GlobalRole
	var categories models.GlobalCategorySpecialist

	response.CustomGormModel = data.CustomGormModel
	response.Avatar = data.Avatar
	response.Fullname = data.Fullname
	response.Email = strings.ToLower(data.Email)
	response.Phone = data.Phone
	response.Address = data.Address
	response.Village = data.Village
	response.District = data.District
	response.City = data.City
	response.Province = data.Province
	response.Country = data.Country
	response.ZipCode = data.ZipCode
	response.Status = data.Status
	response.Gender = data.Gender
	response.Code = data.Code
	response.Experience = data.Experience

	if data.RoleID > 0 {
		roles, _ = GetRoleByIDPlain(data.RoleID)
		response.Role = reqres.GlobalIDNameResponse{
			ID:   int(roles.ID),
			Name: roles.Name,
		}
	}

	if data.CategoryID > 0 {
		categories, _ = GetCategorySpecialistByIDPlain(data.CategoryID)
		response.Category = reqres.GlobalIDNameResponse{
			ID:   int(categories.ID),
			Name: categories.Name,
		}
	}

	return response
}

func GetUsers(roleId, categoryID int, createdAtMarginTop, createdAtMarginBottom string, param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.GlobalUser
	where := "deleted_at IS NULL"

	var modelTotal []models.GlobalUser

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	if createdAtMarginTop != "" {
		where += " AND created_at <= '" + createdAtMarginTop + "'"
	}
	if createdAtMarginBottom != "" {
		where += " AND created_at >= '" + createdAtMarginBottom + "'"
	}

	if roleId > 0 {
		where += " AND role_id = " + strconv.Itoa(roleId)
	}

	if categoryID > 0 {
		where += " AND category_id = " + strconv.Itoa(categoryID)
	}

	if param.Custom != "" {
		where += " AND status = " + param.Custom.(string)
	}

	if param.Search != "" {
		where += " AND (fullname ILIKE '%" + param.Search + "%' OR email ILIKE '%" + param.Search + "%' OR code ILIKE '%" + param.Search + "%')"
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.GlobalUserResponse
	for _, item := range responses {
		responseRefined := BuildUserResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetAllUsers(status int) (users []models.GlobalUser, err error) {
	where := "deleted_at IS NULL"

	if status > 0 {
		where += " AND status = " + strconv.Itoa(status)
	}

	err = config.DB.Where(where).Find(&users).Error

	return
}

func GetUserByIDPlain(id int) (response models.GlobalUser, err error) {
	err = config.DB.First(&response, id).Error

	return
}

func UpdateUser(request models.GlobalUser) (response models.GlobalUser, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func GetUserByID(id int) (response reqres.GlobalUserResponse, err error) {
	var data models.GlobalUser
	err = config.DB.First(&data, id).Error

	response = BuildUserResponse(data)

	return
}

func DeleteUser(request models.GlobalUser) (models.GlobalUser, error) {
	err := config.DB.Delete(&request).Error

	return request, err
}
