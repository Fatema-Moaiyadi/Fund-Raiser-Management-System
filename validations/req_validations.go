package validations

import (
	"github.com/fatema-moaiyadi/fund-raiser-system/models"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"net/mail"
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

	if userInfo.Name == "" {
		return systemerrors.ErrInvalidRequestUserNameEmpty
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
		return systemerrors.ErrFundAmountInvalid
	}

	return nil
}
