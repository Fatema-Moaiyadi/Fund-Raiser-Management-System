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
	ErrInvalidUpdateRequest        = errors.New("request fields invalid, should not contain spaces")
	ErrNameFormatInvalid           = errors.New("request fields invalid, first name/ last name should not contain spaces")
	ErrActiveFunds                 = errors.New("user has ongoing fund raiser(s). Please either delete the fund(s) or wait for it to be completed")
	ErrInvalidRequest              = errors.New("request field(s) invalid, should not be blank")
	ErrInvalidDonationRequest      = errors.New("cannot donate in fund raised by yourself")
	ErrNoUsers                     = errors.New("no users exist in system")
	ErrInvalidFilterRequest        = errors.New("invalid filter parameter")
	ErrNoActiveFunds               = errors.New("no active funds exist in system")
	ErrAmountAlreadyRaised         = errors.New("cannot lower the fund amount. Funds have already been raised upto Rs")
)

func ConvertToErrorWithParams(systemErr error, key string) error {
	return errors.New(fmt.Sprintf("%s %s", systemErr, key))
}
