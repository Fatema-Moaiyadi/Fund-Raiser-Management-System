package models

import "time"

type UserInfo struct {
	UserID    int64      `db:"user_id,omitempty" json:"user_id,omitempty"`
	EmailID   string     `db:"email_id,omitempty" json:"email_id,omitempty"`
	FirstName string     `db:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string     `db:"last_name,omitempty" json:"last_name,omitempty"`
	IsAdmin   bool       `db:"is_admin,omitempty" json:"is_admin,omitempty"`
	Password  string     `db:"password,omitempty" json:"password,omitempty"`
	Status    string     `db:"status,omitempty" json:"status,omitempty"`
	CreatedAt *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type UserInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		UserInfo UserInfo `json:"user_info"`
	} `json:"data"`
}

type DeleteUserResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
