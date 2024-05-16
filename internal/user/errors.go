package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordNotCreated = errors.New("password is not created")
	ErrWrongPassword      = errors.New("wrong password")
	ErrNIPAlreadyExists   = errors.New("NIP already exists")
	ErrValidationFailed   = errors.New("validation failed")
)
