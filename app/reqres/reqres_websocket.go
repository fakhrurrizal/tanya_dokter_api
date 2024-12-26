package reqres

import (
	"tanya_dokter_app/app/models"

	validation "github.com/go-ozzo/ozzo-validation"
)

type GlobalChatRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	File       string `json:"file,omitempty"`
	Status     int    `json:"status"`
	Timestamp  int64  `json:"timestamp"`
}

func (request GlobalChatRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.SenderID, validation.Required),
		validation.Field(&request.ReceiverID, validation.Required),
		validation.Field(&request.Message, validation.Required),
	)
}

type GlobalDataMessagesResponse struct {
	models.CustomGormModel
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	File       string `json:"file,omitempty"`
	Timestamp  int64  `json:"timestamp"`
	Status     int    `json:"status"`
}
