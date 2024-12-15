package repository

import (
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"tanya_dokter_app/config"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

func CreateCategorySpecialist(data *reqres.GlobalCategorySpecialistRequest) (response models.GlobalCategorySpecialist, err error) {

	response = models.GlobalCategorySpecialist{
		Name:        data.Name,
		Description: data.Description,
		Status:      data.Status,
		Image:       data.Image,
		Code:        data.Code,
	}
	var created bool
	for !created {
		err = config.DB.Create(&response).Error
		if err != nil {
			if !config.LoadConfig().EnableIDDuplicationHandling {
				return
			}
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code != "23505" {
					return
				}
			}
		} else {
			created = true
		}
	}

	return
}

func GetCategorySpecialists(param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.GlobalCategorySpecialist
	where := "deleted_at IS NULL"

	var modelTotal []models.GlobalCategorySpecialist

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	if param.Custom != "" {
		where += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		where += " AND (name ILIKE '%" + param.Search + "%' OR description ILIKE '%" + param.Search + "%')"
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.GlobalCategorySpecialistResponse
	for _, item := range responses {
		responseRefined := BuildCategorySpecialistResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetCategorySpecialistByID(id int) (responseRefined reqres.GlobalCategorySpecialistResponse, err error) {
	var response models.GlobalCategorySpecialist
	err = config.DB.First(&response, id).Error

	responseRefined = BuildCategorySpecialistResponse(response)

	return
}

func GetCategorySpecialistByIDPlain(id int) (response models.GlobalCategorySpecialist, err error) {
	err = config.DB.First(&response, id).Error

	return
}

func UpdateCategorySpecialist(request models.GlobalCategorySpecialist) (response models.GlobalCategorySpecialist, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func DeleteCategorySpecialist(request models.GlobalCategorySpecialist) (models.GlobalCategorySpecialist, error) {
	err := config.DB.Delete(&request).Error
	return request, err
}

func BuildCategorySpecialistResponse(data models.GlobalCategorySpecialist) (response reqres.GlobalCategorySpecialistResponse) {

	response.CustomGormModel = data.CustomGormModel
	response.Name = data.Name
	response.Status = data.Status
	response.Description = data.Description
	response.Code = data.Code
	response.Image = data.Image
	response.CreatedAt = data.CreatedAt


	return response
}
