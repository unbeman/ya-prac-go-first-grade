package apperrors

import (
	"errors"
)

var (
	ErrAlreadyExists             = errors.New("already exists")
	ErrInvalidUserCredentials    = errors.New("invalid login or password")
	ErrInvalidToken              = errors.New("invalid token")
	ErrDB                        = errors.New("database error")
	ErrInvalidOrderNumberFormat  = errors.New("invalid order number format")
	ErrNoRecords                 = errors.New("no records")
	ErrNotEnoughPoints           = errors.New("not enough points")
	ErrNoAccrualInfo             = errors.New("no accrual info")
	ErrAccrualServiceUnavailable = errors.New("accrual service is unavailable")
	ErrAccrualConnection         = errors.New("accrual connection error")
)
