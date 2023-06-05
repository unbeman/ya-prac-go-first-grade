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

//type RetryError struct {
//	err      error
//	duration time.Duration
//}
//
//func NewRetryError(duration time.Duration) *RetryError {
//	return &RetryError{err: ErrMustRetry, duration: duration}
//}
//
//func (re *RetryError) Error() string {
//	return fmt.Sprintf("retry error: %v", re.err)
//}
//
//func (re RetryError) GetDuration() time.Duration {
//	return re.duration
//}
