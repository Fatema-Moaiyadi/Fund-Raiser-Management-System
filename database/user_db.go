package database

import "github.com/jmoiron/sqlx"

type userDB struct {
	database *sqlx.DB
}

type UserDatabase interface {
}

func NewUserDB(db *sqlx.DB) UserDatabase {
	return &userDB{
		database: db,
	}
}
