package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"net/http"
)

type fundsHandler struct {
	fundService service.FundService
}

type FundsHandler interface {
	CreateFund() http.HandlerFunc
}

func NewFundsHandler(fundService service.FundService) FundsHandler {
	return &fundsHandler{
		fundService: fundService,
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
