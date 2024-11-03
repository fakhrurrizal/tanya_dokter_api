package models

type GlobalSignin struct {
	CustomGormModel
	UserID      int     `json:"user_id" gorm:"type: int8"`
	BearerToken string  `json:"bearer_token" gorm:"type: text"`
	IPAddress   string  `json:"ip_address" gorm:"type: varchar(255)"`
	UserAgent   string  `json:"user_agent" gorm:"type: text"`
	HostName    string  `json:"host_name" gorm:"type: varchar(255)"`
	AppID       int     `json:"app_id" gorm:"type: int8"`
	CompanyID   int     `json:"company_id" gorm:"type: int8"`
	Latitude    float64 `json:"latitude" gorm:"type: float8"`
	Longitude   float64 `json:"longitude" gorm:"type:float8"`
	ISP         string  `json:"isp" gomr:"type: varchar(255)"`
	City        string  `json:"city" gomr:"type: varchar(255)"`
}

type GlobalSignout struct {
	CustomGormModel
	UserID int    `json:"user_id" gorm:"type: int8"`
	PIN    string `json:"pin" gorm:"type: varchar(255)"`
}
