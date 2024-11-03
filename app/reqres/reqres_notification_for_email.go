package reqres

import "tanya_dokter_app/app/models"

type NotificationForEmailRequest struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
}

type NotificationForEmailResponse struct {
	models.CustomGormModel
	UserID   int    `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
}
