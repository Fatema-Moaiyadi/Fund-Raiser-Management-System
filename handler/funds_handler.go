package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type fundsHandler struct {
	fundService service.FundService
	userService service.UserService
}

type FundsHandler interface {
	CreateFund() http.HandlerFunc
	DonateInFund() http.HandlerFunc
	GetAllActiveFunds() http.HandlerFunc
}

func NewFundsHandler(fundService service.FundService, userService service.UserService) FundsHandler {
	return &fundsHandler{
		fundService: fundService,
		userService: userService,
	}
}

func (fh *fundsHandler) CreateFund() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		createFundRequest := new(models.CreateFundRequest)
		err := json.NewDecoder(request.Body).Decode(createFundRequest)
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

		fundDetails, err := fh.fundService.CreateFund(createFundRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		createFundResponse := new(models.CreateFundResponse)

		res.WriteHeader(http.StatusOK)
		createFundResponse.Code = 0
		createFundResponse.Data.FundInfo = *fundDetails
		response, err := json.Marshal(createFundResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (fh *fundsHandler) DonateInFund() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		donationRequest := new(models.DonationRequest)

		err := json.NewDecoder(request.Body).Decode(donationRequest)
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

		routeVariables := mux.Vars(request)
		fundID, err := strconv.Atoi(routeVariables["fund_id"])
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		tokenPayload := request.Context().Value("claims")
		payload, ok := tokenPayload.(*service.TokenPayload)
		if !ok {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		userInfo, err := fh.userService.FindUser(constants.EmailColumnName, donationRequest.DonatedByEmailID)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		//check if user is admin or token belongs to requested user
		if !payload.IsAdmin && userInfo.UserID != payload.UserID {
			systemerrors.WriteErrorResponse(res, systemerrors.ErrForbidden)
			return
		}

		donationRequest.DonatedByUserID = userInfo.UserID
		donationRequest.DonatedInFund = int64(fundID)

		fundDetails, err := fh.fundService.Donate(donationRequest)

		if err != nil {
			if strings.Contains(err.Error(), systemerrors.ErrLessAmount.Error()) ||
				strings.Contains(err.Error(), systemerrors.ErrMoreAmount.Error()) {
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
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		createFundResponse := new(models.DonationResponse)

		res.WriteHeader(http.StatusOK)
		createFundResponse.Code = 0
		createFundResponse.Data.FundInfo = *fundDetails
		response, err := json.Marshal(createFundResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}

func (fh *fundsHandler) GetAllActiveFunds() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		activeFunds, err := fh.fundService.GetAllActiveFunds()
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		activeFundsResponse := new(models.ActiveFundDetailsResponse)

		res.WriteHeader(http.StatusOK)
		activeFundsResponse.Code = 0
		activeFundsResponse.Data.FundsInfo = activeFunds
		response, err := json.Marshal(activeFundsResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}
		res.Write(response)
	}
}
