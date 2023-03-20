package service

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/database"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
)

type userService struct {
	userDB       database.UserDatabase
	tokenService TokenService
}

type UserService interface {
	Login(email string, password string) (string, error)
}

func NewUserService(userDB database.UserDatabase, ts TokenService) UserService {
	return &userService{
		userDB:       userDB,
		tokenService: ts,
	}
}

func (us *userService) Login(email string, password string) (string, error) {
	userInfo, err := us.userDB.FindUser(email)
	if err != nil {
		return "", err
	}

	if userInfo.Password != password {
		return "", systemerrors.ErrPasswordIncorrect
	}

	token, err := us.tokenService.GenerateToken(userInfo.UserID, userInfo.IsAdmin)
	if err != nil {
		return "", err
	}

	return token, nil
}
