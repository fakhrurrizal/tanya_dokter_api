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

// Fungsi untuk membuat repository baru
func NewChatRepository(data *reqres.GlobalChatRequest) (response models.GlobalMessages, err error) {
	response = models.GlobalMessages{
		SenderID:   data.SenderID,
		ReceiverID: data.ReceiverID,
		Message:    data.Message,
		File:       data.File,
		Timestamp:  data.Timestamp,
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

func GetAllMessages(user_id string, param reqres.ReqPaging) (data reqres.ResPaging) {
	var responses []models.GlobalMessages
	where := "deleted_at IS NULL"

	var modelTotal []models.GlobalMessages

	type TotalResult struct {
		Total       int64
		LastUpdated time.Time
	}
	var totalResult TotalResult
	config.DB.Model(&modelTotal).Select("COUNT(*) AS total, MAX(updated_at) AS last_updated").Scan(&totalResult)

	if param.Custom != "" {
		where += " AND status = " + param.Custom.(string)
	}

	if user_id != "" {
		where += " AND sender_id = " + user_id
	}

	var totalFiltered int64
	config.DB.Model(&modelTotal).Where(where).Count(&totalFiltered)

	config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Sort + " " + param.Order).Where(where).Find(&responses)

	var responsesRefined []reqres.GlobalDataMessagesResponse
	for _, item := range responses {
		responseRefined := BuildDataMessagesResponse(item)

		responsesRefined = append(responsesRefined, responseRefined)
	}

	data = utils.PopulateResPaging(&param, responsesRefined, totalResult.Total, totalFiltered, null.TimeFrom(totalResult.LastUpdated))

	return
}

func BuildDataMessagesResponse(data models.GlobalMessages) (response reqres.GlobalDataMessagesResponse) {

	response.CustomGormModel = data.CustomGormModel
	response.SenderID = data.SenderID
	response.ReceiverID = data.ReceiverID
	response.Message = data.Message
	response.File = data.File
	response.Timestamp = data.Timestamp
	response.CreatedAt = data.CreatedAt

	return response
}
