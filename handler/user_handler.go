package handler

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	"net/http"
)

type userHandler struct {
	userService service.UserService
}

type UserHandler interface {
	LoginHandler() func(res http.ResponseWriter, req *http.Request)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (userHandler *userHandler) LoginHandler() func(res http.ResponseWriter, req *http.Request) {
	return nil
}
