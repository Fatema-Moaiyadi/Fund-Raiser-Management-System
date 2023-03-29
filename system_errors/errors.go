package systemerrors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidRequestEmailEmpty    = errors.New("email is mandatory")
	ErrInvalidRequestPasswordEmpty = errors.New("password is mandatory")
	ErrInvalidRequestUserNameEmpty = errors.New("user name is mandatory")
	ErrUserNotFound                = errors.New("user does not exist in system")
	ErrFundNotFound                = errors.New("fund does not exist in system")
	ErrPasswordIncorrect           = errors.New("incorrect password")
	ErrMissingAuthorizationHeader  = errors.New("token missing in header under Authorization")
	ErrMalformedToken              = errors.New("token is malformed in header")
	ErrInvalidToken                = errors.New("token invalid/expired")
	ErrForbidden                   = errors.New("user is not allowed to access this route")
	ErrInvalidRequestFundNameEmpty = errors.New("fund name is mandatory")
	ErrAmountInvalid               = errors.New("amount should be greater than 0")
	ErrLessAmount                  = errors.New("amount should be greater than")
	ErrMoreAmount                  = errors.New("amount should be less than")
	ErrFundInactive                = errors.New("requested fund is inactive")
	ErrInvalidRequest              = errors.New("request fields invalid, should not contain spaces")
	ErrNameFormatInvalid           = errors.New("request fields invalid, first name/ last name should not contain spaces")
)

func ConvertToUserSpecificError(systemErr error, err string) error {
	return errors.New(fmt.Sprintf("%s %s", systemErr, err))
}
