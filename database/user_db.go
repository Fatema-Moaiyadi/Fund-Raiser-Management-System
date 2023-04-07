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
	FindUserForLogin(emailID string) (*models.UserInfo, error)
	CreateUser(userDetails *models.UserInfo) error
	UpdateUserByID(userID int64, updateParams map[string]string) (*models.UpdateUser, error)
	DeleteUserByID(userID int64) error
	GetUserInfoByFilters(filterParams map[string]interface{}) (*models.UserDetailedInfo, error)
	GetAllUsersInfo() ([]models.UserDetailedInfo, error)
}

func NewUserDB(db *sqlx.DB) UserDatabase {
	return &userDB{
		database: db,
	}
}

func (ud *userDB) FindUserForLogin(emailID string) (*models.UserInfo, error) {
	userInfo := new(models.UserInfo)
	err := ud.database.Get(userInfo, findUserForLoginQuery, emailID)

	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (ud *userDB) FindUser(filterKey string, filterValue interface{}) (*models.UserInfo, error) {
	userInfo := new(models.UserInfo)

	query := fmt.Sprintf(getUserQuery, filterKey)
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

	result, err := tx.Exec(deleteUserByID, time.Now(), userID)
	affectedRows, e := result.RowsAffected()
	if e != nil {
		return e
	}

	if affectedRows == 0 {
		return systemerrors.ErrUserNotFound
	}

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

func (ud *userDB) GetUserInfoByFilters(filterParams map[string]interface{}) (*models.UserDetailedInfo, error) {
	userDetails := new(models.UserDetailedInfo)

	var (
		whereClause  []string
		filterValues []interface{}
		i            = 1
	)
	for key, value := range filterParams {
		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", key, i))
		i++
		filterValues = append(filterValues, value)
	}

	query := fmt.Sprintf(getUserInfoByIDQuery, strings.Join(whereClause, " AND "))

	err := ud.database.Get(&userDetails.UserInfo, query, filterValues...)
	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	err = ud.database.Select(&userDetails.RaisedFundsInfo, getFundsRaisedByUserIDQuery, userDetails.UserInfo.UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	err = ud.database.Select(&userDetails.DonationsInfo, getDonationsByUserIDQuery, userDetails.UserInfo.UserID)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return userDetails, nil
}

func (ud *userDB) GetAllUsersInfo() ([]models.UserDetailedInfo, error) {
	usersInfo := make([]models.UserInfo, 0)
	err := ud.database.Select(&usersInfo, getAllUserInfoQuery)
	if err == sql.ErrNoRows {
		return nil, systemerrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	allUserDetails := make([]models.UserDetailedInfo, len(usersInfo))
	for i, userInfo := range usersInfo {
		allUserDetails[i].UserInfo = userInfo
		err = ud.database.Select(&allUserDetails[i].RaisedFundsInfo, getFundsRaisedByUserIDQuery, userInfo.UserID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}

		err = ud.database.Select(&allUserDetails[i].DonationsInfo, getDonationsByUserIDQuery, userInfo.UserID)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
		}
	}

	return allUserDetails, nil
}
