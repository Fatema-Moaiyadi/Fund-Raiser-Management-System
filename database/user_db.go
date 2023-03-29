package database

import (
	"database/sql"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type userDB struct {
	database *sqlx.DB
}

type UserDatabase interface {
	FindUser(filterKey string, filterValue interface{}) (*models.UserInfo, error)
	CreateUser(userDetails *models.UserInfo) error
	UpdateUserByID(userID int64, updateParams map[string]string) (*models.UpdateUser, error)
	DeleteUserByID(userID int64) error
}

func NewUserDB(db *sqlx.DB) UserDatabase {
	return &userDB{
		database: db,
	}
}

func (ud *userDB) FindUser(filterKey string, filterValue interface{}) (*models.UserInfo, error) {
	userInfo := new(models.UserInfo)
	query := fmt.Sprintf(findUserQuery, filterKey)
	err := ud.database.Get(userInfo, query, filterValue)

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
		userDetails.FirstName,
		userDetails.LastName,
		userDetails.Password,
		userDetails.IsAdmin,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (ud *userDB) UpdateUserByID(userID int64, updateParams map[string]string) (*models.UpdateUser, error) {
	var (
		setValues []interface{}
		setKeys   []string
		getValues []string
		index     = 1
	)

	for key, value := range updateParams {
		setKeys = append(setKeys, key+"=$"+fmt.Sprintf("%d", index))
		setValues = append(setValues, value)
		getValues = append(getValues, key)
		index++
	}

	tx, err := ud.database.Begin()
	if err != nil {
		return nil, err
	}

	updateQuery := updateUserSetClause + strings.Join(setKeys, ",") + fmt.Sprintf(whereUserIDClause, index, index+1)

	result, err := tx.Exec(updateQuery, append(setValues, time.Now(), userID)...)

	affectedRows, e := result.RowsAffected()
	if e != nil {
		return nil, e
	}

	if affectedRows == 0 {
		return nil, systemerrors.ErrUserNotFound
	}
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return nil, e
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedDetails := new(models.UpdateUser)
	getQuery := fmt.Sprintf(getUserByIDWithUpdatedInfo, strings.Join(getValues, ","))
	err = ud.database.Get(updatedDetails, getQuery, userID)
	if err != nil {
		return nil, err
	}

	return updatedDetails, nil
}

func (ud *userDB) DeleteUserByID(userID int64) error {
	tx, err := ud.database.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(deleteUserByID, time.Now(), userID)
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return e
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
