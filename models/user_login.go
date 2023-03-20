package models

import (
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"net/mail"
)

type UserLoginRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Code int `json:"code"`
	Data struct {
		AuthToken string `json:"auth_token"`
	} `json:"data"`
}

func (req UserLoginRequest) Validate() error {
	_, err := mail.ParseAddress(req.EmailID)
	if err != nil {
		return systemerrors.ErrInvalidLoginRequestEmailEmpty
	}

	if req.Password == "" {
		return systemerrors.ErrInvalidLoginRequestPasswordEmpty
	}

	return nil
}
