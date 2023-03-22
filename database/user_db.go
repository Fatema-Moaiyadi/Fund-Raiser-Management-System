package database

import (
	"database/sql"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type userDB struct {
	database *sqlx.DB
}

type UserDatabase interface {
	FindUser(email string) (*models.UserInfo, error)
	CreateUser(userDetails *models.UserInfo) error
}

func NewUserDB(db *sqlx.DB) UserDatabase {
	return &userDB{
		database: db,
	}
}

func (ud *userDB) FindUser(email string) (*models.UserInfo, error) {
	userInfo := new(models.UserInfo)
	err := ud.database.Get(userInfo, findUserQuery, email)

	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (ud *userDB) CreateUser(userDetails *models.UserInfo) error {
	err := ud.database.Get(userDetails, insertUserQuery,
		userDetails.EmailID,
		userDetails.Name,
		userDetails.Password,
		userDetails.IsAdmin,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}
