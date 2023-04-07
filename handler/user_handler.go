package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/fatema-moaiyadi/fund-raiser-system/validations"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type userHandler struct {
	userService service.UserService
}

type UserHandler interface {
	LoginHandler() http.HandlerFunc
	CreateUser() http.HandlerFunc
	UpdateUserInfo() http.HandlerFunc
	DeleteUserByID() http.HandlerFunc
	GetUserInfoByID() http.HandlerFunc
	GetAllUsersInfo() http.HandlerFunc
	GetUserInfoByFilters() http.HandlerFunc
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
			systemerrors.WriteErrorResponse(res, err)
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

		userID, err := strconv.Atoi(mux.Vars(request)["fund_id"])
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		err = userHandler.userService.DeleteUserByID(int64(userID))
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		deleteUserResponse := new(models.DeleteUserResponse)
		res.WriteHeader(http.StatusOK)
		deleteUserResponse.Code = 0
		deleteUserResponse.Message = fmt.Sprintf("User with user id %d deleted successfully", userID)

		response, err := json.Marshal(deleteUserResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (userHandler *userHandler) GetUserInfoByID() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		requestParams := mux.Vars(request)
		userID, err := strconv.Atoi(requestParams["user_id"])
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		tokenPayload := request.Context().Value("claims")
		payload := tokenPayload.(*service.TokenPayload)

		if payload.UserID != int64(userID) {
			systemerrors.WriteErrorResponse(res, systemerrors.ErrForbidden)
			return
		}

		filterParams := make(map[string]interface{})
		filterParams[constants.UserIDColumnName] = int64(userID)
		userDetails, err := userHandler.userService.GetUserInfoByFilters(filterParams)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		userDetailsResponse := new(models.GetUserDetailsByIDResponse)
		res.WriteHeader(http.StatusOK)
		userDetailsResponse.Code = 0
		userDetailsResponse.Data.UserDetails = *userDetails
		response, err := json.Marshal(userDetailsResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (userHandler *userHandler) GetAllUsersInfo() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		allUserDetails, err := userHandler.userService.GetAllUsersInfo()
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		allUserDetailsResponse := new(models.GetAllUserDetailsResponse)
		res.WriteHeader(http.StatusOK)
		allUserDetailsResponse.Code = 0
		allUserDetailsResponse.Data.AllUsersInfo = allUserDetails
		response, err := json.Marshal(allUserDetailsResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (userHandler *userHandler) GetUserInfoByFilters() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		filterParams := make(map[string]interface{})

		emailID := request.URL.Query().Get(constants.EmailColumnName)
		if emailID != "" {
			filterParams[constants.EmailColumnName] = emailID
		}

		firstName := request.URL.Query().Get(constants.FirstNameColumnName)
		if firstName != "" {
			filterParams[constants.FirstNameColumnName] = firstName
		}

		lastName := request.URL.Query().Get(constants.LastNameColumnName)
		if lastName != "" {
			filterParams[constants.LastNameColumnName] = lastName
		}

		userDetails, err := userHandler.userService.GetUserInfoByFilters(filterParams)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		userDetailsResponse := new(models.GetUserDetailsByIDResponse)
		res.WriteHeader(http.StatusOK)
		userDetailsResponse.Code = 0
		userDetailsResponse.Data.UserDetails = *userDetails
		response, err := json.Marshal(userDetailsResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}
