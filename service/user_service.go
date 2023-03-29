package service

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
	"github.com/fatema-moaiyadi/fund-raiser-system/database"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/fatema-moaiyadi/fund-raiser-system/validations"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userDB       database.UserDatabase
	tokenService TokenService
}

type UserService interface {
	Login(email string, password string) (string, error)
	CreateUser(userDetails *models.UserInfo) error
	FindUser(filterKey string, filterValue interface{}) (*models.UserInfo, error)
	UpdateUserByID(userID int64, updateUserReq *models.UpdateUser) (*models.UpdateUser, error)
}

func NewUserService(userDB database.UserDatabase, ts TokenService) UserService {
	return &userService{
		userDB:       userDB,
		tokenService: ts,
	}
}

func (us *userService) Login(email string, password string) (string, error) {
	userInfo, err := us.userDB.FindUser(constants.EmailColumnName, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password))
	if err != nil {
		return "", systemerrors.ErrPasswordIncorrect
	}

	token, err := us.tokenService.GenerateToken(userInfo.UserID, userInfo.IsAdmin)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *userService) CreateUser(userDetails *models.UserInfo) error {
	err := validations.ValidateCreateUserReq(userDetails)
	if err != nil {
		return err
	}

	randomPassword, err := password.Generate(10, 0, 0, false, true)
	if err != nil {
		return err
	}

	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(randomPassword), 14)
	if err != nil {
		return err
	}

	userDetails.Password = string(hashedPasswordByte)

	err = us.userDB.CreateUser(userDetails)
	if err != nil {
		return err
	}

	//need to show user the actual un-hashed or decrypted password
	//for first time login
	userDetails.Password = randomPassword

	return nil
}

func (us *userService) FindUser(filterKey string, filterValue interface{}) (*models.UserInfo, error) {
	userInfo, err := us.userDB.FindUser(filterKey, filterValue)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (us *userService) UpdateUserByID(userID int64, updateUserReq *models.UpdateUser) (*models.UpdateUser, error) {
	err := validations.ValidateUpdateUserRequest(updateUserReq)
	if err != nil {
		return nil, err
	}
	updateParams := make(map[string]string)

	if updateUserReq.FirstName != "" {
		updateParams[constants.FirstNameColumnName] = updateUserReq.FirstName
	}

	if updateUserReq.LastName != "" {
		updateParams[constants.LastNameColumnName] = updateUserReq.LastName
	}

	if updateUserReq.Password != "" {
		hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(updateUserReq.Password), 14)
		if err != nil {
			return nil, err
		}
		updateParams[constants.PasswordColumnName] = string(hashedPasswordByte)
	}

	updatedInfo, err := us.userDB.UpdateUserByID(userID, updateParams)
	if err != nil {
		return nil, err
	}

	if updatedInfo.Password != "" {
		updatedInfo.Password = updateUserReq.Password
	}

	return updatedInfo, nil
}
