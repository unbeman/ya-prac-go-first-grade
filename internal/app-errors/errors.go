package app_errors

import (
	"errors"
)

var (
	ErrAlreadyExists          = errors.New("already exists")
	ErrInvalidUserCredentials = errors.New("invalid login or password")
	ErrInvalidToken           = errors.New("invalid token")
	ErrDb                     = errors.New("database error")
	//ErrInvalidContentType       = errors.New("invalid content type")
	ErrInvalidOrderNumberFormat = errors.New("invalid order number format")
	ErrNoRecords                = errors.New("no records")
	ErrNotEnoughPoints          = errors.New("not enough points")
	ErrNoAccrualInfo            = errors.New("no accrual info")
)
