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
	UpdateFund() http.HandlerFunc
	DeleteFund() http.HandlerFunc
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
			systemerrors.WriteErrorResponse(res, err)
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
			systemerrors.WriteErrorResponse(res, err)
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

func (fh *fundsHandler) UpdateFund() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		updateFundRequest := new(models.UpdateFund)

		err := json.NewDecoder(request.Body).Decode(&updateFundRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		fundID, err := strconv.Atoi(mux.Vars(request)["fund_id"])
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		updatedDetails, err := fh.fundService.UpdateFundByID(int64(fundID), updateFundRequest)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		updatedFundResponse := new(models.UpdateFundResponse)

		updatedFundResponse.Code = 0
		updatedFundResponse.Message = "Fund updated successfully"
		updatedFundResponse.Data.UpdatedInfo = *updatedDetails

		response, err := json.Marshal(updatedFundResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}

		res.Write(response)
	}
}

func (fh *fundsHandler) DeleteFund() http.HandlerFunc {
	return func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		fundID, err := strconv.Atoi(mux.Vars(request)["fund_id"])
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		err = fh.fundService.DeleteFundByID(int64(fundID))
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		deleteFundResponse := new(models.DeleteFundResponse)

		deleteFundResponse.Code = 0
		deleteFundResponse.Message = fmt.Sprintf("Fund with fund id %d deleted successfully", fundID)
		response, err := json.Marshal(deleteFundResponse)
		if err != nil {
			fmt.Fprintf(res, "Decoding error")
			return
		}

		res.Write(response)
	}
}
