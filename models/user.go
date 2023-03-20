package models

type UserInfo struct {
	UserID   int    `db:"user_id"`
	EmailID  string `db:"email_id"`
	IsAdmin  bool   `db:"is_admin"`
	Password string `db:"password"`
}
