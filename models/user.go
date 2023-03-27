package models

import "time"

type UserInfo struct {
	UserID    int64     `db:"user_id,omitempty" json:"user_id,omitempty"`
	EmailID   string    `db:"email_id,omitempty" json:"email_id,omitempty"`
	Name      string    `db:"name,omitempty" json:"name,omitempty"`
	IsAdmin   bool      `db:"is_admin,omitempty" json:"is_admin,omitempty"`
	Password  string    `db:"password,omitempty" json:"password,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type UserInfoResponse struct {
	Code int `json:"code"`
	Data struct {
		UserInfo UserInfo `json:"user_info"`
	} `json:"data"`
}
