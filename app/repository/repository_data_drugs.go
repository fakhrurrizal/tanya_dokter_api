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

func CreateDataDrugs(data *reqres.GlobalDataDrugsRequest) (response models.GlobalDataDrugs, err error) {

	response = models.GlobalDataDrugs{
		Name:        data.Name,
		Description: data.Description,
		Image:       data.Image,
		Code:        data.Code,
		Usage:       data.Usage,
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

func GetDataDrugs(param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.GlobalDataDrugs
	where := "deleted_at IS NULL"

	var modelTotal []models.GlobalDataDrugs

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

	var responsesRefined []reqres.GlobalDataDrugsResponse
	for _, item := range responses {
		responseRefined := BuildDataDrugsResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func GetDataDrugsByID(id int) (responseRefined reqres.GlobalDataDrugsResponse, err error) {
	var response models.GlobalDataDrugs
	err = config.DB.First(&response, id).Error

	responseRefined = BuildDataDrugsResponse(response)

	return
}

func GetDataDrugsByIDPlain(id int) (response models.GlobalDataDrugs, err error) {
	err = config.DB.First(&response, id).Error

	return
}

func UpdateDataDrugs(request models.GlobalDataDrugs) (response models.GlobalDataDrugs, err error) {
	err = config.DB.Save(&request).Scan(&response).Error

	return
}

func DeleteDataDrugs(request models.GlobalDataDrugs) (models.GlobalDataDrugs, error) {
	err := config.DB.Delete(&request).Error
	return request, err
}

func BuildDataDrugsResponse(data models.GlobalDataDrugs) (response reqres.GlobalDataDrugsResponse) {

	response.CustomGormModel = data.CustomGormModel
	response.Name = data.Name
	response.Description = data.Description
	response.Code = data.Code
	response.Image = data.Image
	response.Usage = data.Usage
	response.CreatedAt = data.CreatedAt

	return response
}
