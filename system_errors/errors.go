package systemerrors

import "errors"

var (
	ErrInvalidLoginRequestEmailEmpty    = errors.New("email is mandatory")
	ErrInvalidLoginRequestPasswordEmpty = errors.New("password is mandatory")
	ErrUserNotFound                     = errors.New("user does not exist in system")
	ErrPasswordIncorrect                = errors.New("incorrect password")
)
