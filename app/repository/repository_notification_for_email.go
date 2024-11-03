package repository

import (
	"strconv"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"tanya_dokter_app/config"
	"time"

	"github.com/guregu/null"
)

func CreateNotificationForEmail(data *reqres.NotificationForEmailRequest) (response models.NotificationForEmail, err error) {
	response = models.NotificationForEmail{
		UserID:   data.UserID,
		FullName: data.FullName,
		Email:    data.Email,
		Status:   data.Status,
	}

	err = config.DB.Create(&response).Error
	if err != nil {
		return
	}

	return
}

func BuildNotificationForEmailResponse(data models.NotificationForEmail) (response reqres.NotificationForEmailResponse) {
	response = reqres.NotificationForEmailResponse{
		CustomGormModel: data.CustomGormModel,
		UserID:          data.UserID,
		FullName:        data.FullName,
		Email:           data.Email,
		Status:          data.Status,
	}
	return
}

func GetNotificationForEmails(userID int, param reqres.ReqPaging) (data reqres.ResPaging) {
	var out []models.NotificationForEmail
	where := "deleted_at IS NULL"
	if userID != 0 {
		where += " AND user_id = " + strconv.Itoa(userID)
	}
	if param.Custom != "" {
		where += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		where += " AND (full_name ILIKE '%" + param.Search + "%' OR email ILIKE '%" + param.Search + "%')"
	}

	var modelTotal []models.NotificationForEmail

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&out)

	var responses []reqres.NotificationForEmailResponse
	for _, item := range out {
		response := BuildNotificationForEmailResponse(item)
		responses = append(responses, response)
	}

	data = utils.PopulateResPaging(&param, responses, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetAllNotificationForEmails() (data []models.NotificationForEmail, err error) {
	err = config.DB.Where("status = true").Find(&data).Error
	return
}

func GetNotificationForEmailByID(id int) (response reqres.NotificationForEmailResponse, err error) {
	var out models.NotificationForEmail
	err = config.DB.First(&out, id).Error
	response = BuildNotificationForEmailResponse(out)
	return
}

func GetNotificationForEmailByIDPlain(id int) (response models.NotificationForEmail, err error) {
	err = config.DB.First(&response, id).Error
	return
}

func UpdateNotificationForEmailByID(data models.NotificationForEmail) (response models.NotificationForEmail, err error) {
	err = config.DB.Save(&data).Scan(&response).Error
	return
}

func DeleteNotificationForEmailByID(data models.NotificationForEmail) (response models.NotificationForEmail, err error) {
	err = config.DB.Delete(&data).Scan(&response).Error
	return
}
