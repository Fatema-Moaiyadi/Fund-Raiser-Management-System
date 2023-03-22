package database

const (
	findUserQuery   = "SELECT * from users where email_id = $1"
	insertUserQuery = "INSERT INTO users (email_id,name,password,is_admin,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id,email_id,name,password,is_admin,created_at,updated_at"

	insertFundQuery = "INSERT INTO funds (raised_by_user_id,name,amount,status,created_at,updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING fund_id,raised_by_user_id,name,amount,status,created_at,updated_at"
)
