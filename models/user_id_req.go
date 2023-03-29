package models

type UserIDRequest struct {
	UserID int64 `json:"user_id"`
}

type DeleteUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
