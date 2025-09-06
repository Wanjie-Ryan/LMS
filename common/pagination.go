package common

type Pagination struct{
	Limit int `query:"limit" json:"limit"`
	Page int `query:"page" json:"page"`
	Sort string `query:"sort" json:"sort"`
	TotalRows int64 `json:"total_rows"`
	TotalPage int `json:"total_page"`
	Data interface{} `json:"data"`
}