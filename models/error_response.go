package models

type ErrorResponse struct {
	Code  int `json:"code"`
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}
