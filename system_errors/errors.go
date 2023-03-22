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
	ErrPasswordIncorrect           = errors.New("incorrect password")
	ErrMissingAuthorizationHeader  = errors.New("token missing in header under Authorization")
	ErrMalformedToken              = errors.New("token is malformed in header")
	ErrInvalidToken                = errors.New("token invalid/expired")
	ErrForbidden                   = errors.New("user is not allowed to access this route")
	ErrInvalidRequestFundNameEmpty = errors.New("fund name is mandatory")
	ErrFundAmountInvalid           = errors.New("fund amount should be greater than 0")
)

func ConvertToCustomError(customErr error, actualError string) error {
	return errors.New(fmt.Sprintf("%s - %s", customErr, actualError))
}
