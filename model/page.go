package model

type Page struct {
	Current   int `json:"current"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}
