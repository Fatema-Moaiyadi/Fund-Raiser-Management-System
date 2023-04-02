package systemerrors

import (
	"encoding/json"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	"net/http"
)

func WriteErrorResponse(res http.ResponseWriter, err error) {
	var statusCode int
	switch err {
	case ErrMissingAuthorizationHeader, ErrMalformedToken:
		statusCode = http.StatusBadRequest
	case ErrInvalidToken, ErrPasswordIncorrect:
		statusCode = http.StatusUnauthorized
	case ErrForbidden:
		statusCode = http.StatusForbidden
	case ErrInvalidRequestUserNameEmpty, ErrInvalidRequestEmailEmpty,
		ErrInvalidRequestFundNameEmpty, ErrAmountInvalid, ErrInvalidUpdateRequest,
		ErrNameFormatInvalid, ErrInvalidRequest, ErrInvalidDonationRequest:
		statusCode = http.StatusBadRequest
	case ErrUserNotFound, ErrFundNotFound, ErrNoUsers:
		statusCode = http.StatusNotFound
	case ErrActiveFunds:
		statusCode = http.StatusMethodNotAllowed
	default:
		statusCode = http.StatusInternalServerError
	}

	res.Header().Set("Content-Type", "application/json")

	res.WriteHeader(statusCode)
	errorRes := new(models.ErrorResponse)
	errorRes.Error.Message = err.Error()
	errorRes.Error.Status = statusCode
	errorRes.Code = -1

	response, err := json.Marshal(errorRes)
	if err != nil {
		fmt.Fprintf(res, "Decoding error")
	}
	res.Write(response)
}
