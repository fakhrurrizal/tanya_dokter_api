package models

type GlobalMessages struct {
	CustomGormModel
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Message    string `json:"message"`
	File       string `json:"file,omitempty"`
	Timestamp  int64  `json:"timestamp"`
}
