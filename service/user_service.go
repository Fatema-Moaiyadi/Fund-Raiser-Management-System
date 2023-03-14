package service

import "github.com/fatema-moaiyadi/fund-raiser-system/database"

type userService struct {
	userDB database.UserDatabase
}

type UserService interface {
}

func NewUserService(userDB database.UserDatabase) UserService {
	return &userService{
		userDB: userDB,
	}
}
