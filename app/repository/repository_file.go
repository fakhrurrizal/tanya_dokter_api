package repository

import (
	"strconv"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/app/utils"
	"tanya_dokter_app/config"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
)

func SaveFile(data *models.GlobalFile) (err error) {
	var created bool
	for !created {
		err = config.DB.Create(&data).Error
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

func GetFile(id int, param reqres.ReqPaging) (data reqres.ResPaging, err error) {
	search := "user_id = " + strconv.Itoa(id)
	if param.Search != "" {
		search = " AND filename ILIKE '%" + param.Search + "%' "
	}
	token := param.Custom.(string)
	if token != "" {
		search += " AND token ILIKE '%" + token + "%' "
	}
	var files []models.GlobalFile
	err = config.DB.Where(search).Order(param.Sort + " " + param.Order).Limit(param.Limit).Offset(param.Offset).Find(&files).Error
	if err != nil {
		return
	}

	var modelTotal []models.GlobalFile

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	var totalFiltered int64

	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	config.DB.Model(&modelTotal).Where(search).Count(&totalFiltered)

	for i, file := range files {
		files[i].FullUrl = config.LoadConfig().BaseUrl + "/assets/uploads/" + file.Filename

		if config.LoadConfig().EnableEncodeID {
			if file.EncodedID == "" {
				files[i].EncodedID = utils.EndcodeID(int(file.ID))
			}
		}
	}
	data = utils.PopulateResPaging(&param, files, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))
	return
}

func GetFileByToken(token string, id int64, companyID int) (data models.GlobalFile, err error) {
	err = config.DB.Where("user_id = ? and token = ?", id, token).First(&data).Error

	return
}
