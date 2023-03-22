package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/fatema-moaiyadi/fund-raiser-system/validations"
	"net/http"
)

type userHandler struct {
	userService service.UserService
}

type UserHandler interface {
	LoginHandler() http.HandlerFunc
	CreateUser() http.HandlerFunc
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (userHandler *userHandler) LoginHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := new(models.UserLoginRequest)

		res.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(req.Body).Decode(loginReq)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			errorRes := new(models.ErrorResponse)
			errorRes.Error.Message = "Invalid Request"
			errorRes.Error.Status = http.StatusBadRequest
			errorRes.Code = -1

			response, err := json.Marshal(errorRes)
			if err != nil {
				fmt.Fprintf(res, "Decoding error")
				return
			}
			res.Write(response)
			return
		}

		err = validations.ValidateLoginRequest(*loginReq)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			errorRes := new(models.ErrorResponse)
			errorRes.Error.Message = err.Error()
			errorRes.Error.Status = http.StatusBadRequest
			errorRes.Code = -1

			response, err := json.Marshal(errorRes)
			if err != nil {
				fmt.Fprintf(res, "Decoding error")
			}
			res.Write(response)
			return
		}

		token, err := userHandler.userService.Login(loginReq.EmailID, loginReq.Password)
		if err != nil {
			errorRes := new(models.ErrorResponse)
			if err == systemerrors.ErrUserNotFound {
				res.WriteHeader(http.StatusNotFound)
				errorRes.Error.Status = http.StatusNotFound
			} else if err == systemerrors.ErrPasswordIncorrect {
				res.WriteHeader(http.StatusUnauthorized)
				errorRes.Error.Status = http.StatusUnauthorized
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				errorRes.Error.Status = http.StatusInternalServerError
			}
			errorRes.Error.Message = err.Error()
			errorRes.Code = -1

			response, err := json.Marshal(errorRes)
			if err != nil {
				fmt.Fprintf(res, "Decoding error")
				return
			}
			res.Write(response)
			return
		}

		res.WriteHeader(http.StatusOK)

		loginResp := new(models.UserLoginResponse)
		loginResp.Code = 0
		loginResp.Data.AuthToken = token

		response, err := json.Marshal(loginResp)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (userHandler *userHandler) CreateUser() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		createUserRequest := new(models.UserInfo)

		err := json.NewDecoder(request.Body).Decode(createUserRequest)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			errorRes := new(models.ErrorResponse)
			errorRes.Error.Message = "Invalid Request"
			errorRes.Error.Status = http.StatusBadRequest
			errorRes.Code = -1

			response, err := json.Marshal(errorRes)
			if err != nil {
				fmt.Fprintf(res, "Decoding error")
			}
			res.Write(response)
			return
		}

		err = userHandler.userService.CreateUser(createUserRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		createUserResponse := new(models.UserInfoResponse)

		res.WriteHeader(http.StatusOK)
		createUserResponse.Code = 0
		createUserResponse.Data.UserInfo = *createUserRequest
		response, err := json.Marshal(createUserResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}
