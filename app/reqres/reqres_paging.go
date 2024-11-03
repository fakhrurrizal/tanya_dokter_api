package reqres

import "github.com/guregu/null"

type ReqPaging struct {
	Page   int         `default:"1"`
	Search string      `default:""`
	Limit  int         `default:"10"`
	Offset int         `default:"0"`
	Sort   string      `default:"ASC"`
	Order  string      `default:"id"`
	Custom interface{} `default:""`
}

type ResPaging struct {
	Draw            int         `json:"draw"`
	TotalData       int         `json:"recordsTotal"`
	RecordsFiltered int         `json:"recordsFiltered"`
	Error           string      `json:"error"`
	Status          int         `default:"200" json:"status"`
	Messages        string      `default:"Success" json:"message"`
	Data            interface{} `default:"[]" json:"data"`
	Search          string      `default:"" json:"search"`
	Next            bool        `default:"false" json:"next"`
	Back            bool        `default:"false" json:"back"`
	Limit           int         `default:"10" json:"limit"`
	Offset          int         `default:"0" json:"offset"`
	TotalPage       int         `default:"0" json:"total_page"`
	CurrentPage     int         `default:"1" json:"current_page"`
	Sort            string      `default:"ASC" json:"sort"`
	Order           string      `default:"id" json:"order"`
	Summary         interface{} `json:"summary,omitempty"`
	LastUpdated     null.Time   `json:"last_updated"`
}
