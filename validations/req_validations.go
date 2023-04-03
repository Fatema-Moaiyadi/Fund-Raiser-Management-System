package validations

import (
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/constants"
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"net/mail"
	"strings"
)

func ValidateLoginRequest(req models.UserLoginRequest) error {
	_, err := mail.ParseAddress(req.EmailID)
	if err != nil {
		return systemerrors.ErrInvalidRequestEmailEmpty
	}

	if req.Password == "" {
		return systemerrors.ErrInvalidRequestPasswordEmpty
	}

	return nil
}

func ValidateCreateUserReq(userInfo *models.UserInfo) error {
	_, err := mail.ParseAddress(userInfo.EmailID)
	if err != nil {
		return systemerrors.ErrInvalidRequestEmailEmpty
	}

	if userInfo.FirstName == "" {
		return systemerrors.ErrInvalidRequestUserNameEmpty
	}

	if strings.Contains(userInfo.FirstName, " ") ||
		strings.Contains(userInfo.LastName, " ") {
		return systemerrors.ErrNameFormatInvalid
	}

	return nil
}

func ValidateCreateFundReq(fundReq *models.CreateFundRequest) error {
	if fundReq.FundName == "" {
		return systemerrors.ErrInvalidRequestFundNameEmpty
	}

	if fundReq.RaisedByUserEmail == "" {
		return systemerrors.ErrInvalidRequestEmailEmpty
	}

	if fundReq.TotalAmount <= 0 {
		return systemerrors.ErrAmountInvalid
	}

	return nil
}

func ValidateDonateRequest(donationRequest models.DonationRequest, amountRaised int64, fundDetails *models.FundDetails) error {
	totalFundAmount := fundDetails.TotalAmount

	if donationRequest.DonationAmount <= 0 {
		return systemerrors.ErrAmountInvalid
	}

	if donationRequest.DonationAmount < constants.MinDonationAmount {
		return systemerrors.ConvertToUserSpecificError(systemerrors.ErrLessAmount, fmt.Sprintf("%d", constants.MinDonationAmount))
	}

	if donationRequest.DonationAmount > totalFundAmount-amountRaised {
		return systemerrors.ConvertToUserSpecificError(systemerrors.ErrMoreAmount, fmt.Sprintf("%d", totalFundAmount-amountRaised))
	}

	if donationRequest.DonationAmount > constants.MaxDonationAmount {
		return systemerrors.ConvertToUserSpecificError(systemerrors.ErrMoreAmount, fmt.Sprintf("%d", constants.MaxDonationAmount))
	}

	if fundDetails.FundStatus != models.IN_PROGRESS.String() {
		return systemerrors.ErrFundInactive
	}

	if fundDetails.RaisedByUserID == donationRequest.DonatedByUserID {

	}
	return nil
}

func ValidateUpdateUserRequest(updateRequest *models.UpdateUser) error {
	//should not contain spaces to avoid sql injection
	if strings.Contains(updateRequest.FirstName, " ") ||
		strings.Contains(updateRequest.LastName, " ") ||
		strings.Contains(updateRequest.Password, " ") {
		return systemerrors.ErrInvalidUpdateRequest
	}

	return nil
}

func ValidateUserIDInRequest(request *models.UserIDRequest) error {
	if request.UserID == 0 {
		return systemerrors.ErrInvalidRequest
	}
	return nil
}

func ValidateFilterKeysInRequest(filters map[string]interface{}) error {
	filtersValid := true
	for key, _ := range filters {
		if key == constants.EmailColumnName || key == constants.UserIDColumnName {
			continue
		} else if key == constants.FirstNameColumnName {
			_, ok := filters[constants.LastNameColumnName]
			if !ok {
				filtersValid = false
				break
			}
		} else if key == constants.LastNameColumnName {
			_, ok := filters[constants.FirstNameColumnName]
			if !ok {
				filtersValid = false
				break
			}
		} else {
			filtersValid = false
			break
		}
	}

	if !filtersValid {
		return systemerrors.ErrInvalidFilterRequest
	}
	return nil
}
