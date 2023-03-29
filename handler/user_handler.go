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
	UpdateUserInfo() http.HandlerFunc
	DeleteUserByID() http.HandlerFunc
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
			systemerrors.WriteErrorResponse(res, err)
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
			systemerrors.WriteErrorResponse(res, err)
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

func (userHandler *userHandler) UpdateUserInfo() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		updateUserRequest := new(models.UpdateUser)

		err := json.NewDecoder(request.Body).Decode(updateUserRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		tokenPayload := request.Context().Value("claims").(*service.TokenPayload)

		userID := tokenPayload.UserID

		updatedInfo, err := userHandler.userService.UpdateUserByID(userID, updateUserRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		updateUserResponse := new(models.UpdateUserResponse)
		res.WriteHeader(http.StatusOK)
		updateUserResponse.Code = 0
		updateUserResponse.Message = "User updated successfully"
		updateUserResponse.Data.UpdatedInfo = *updatedInfo
		response, err := json.Marshal(updateUserResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (userHandler *userHandler) DeleteUserByID() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		deleteUserRequest := new(models.UserIDRequest)

		err := json.NewDecoder(request.Body).Decode(deleteUserRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		err = userHandler.userService.DeleteUserByID(deleteUserRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		deleteUserResponse := new(models.DeleteUserResponse)
		res.WriteHeader(http.StatusOK)
		deleteUserResponse.Code = 0
		deleteUserResponse.Message = fmt.Sprintf("User with user id %d deleted successfully", deleteUserRequest.UserID)

		response, err := json.Marshal(deleteUserResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}
