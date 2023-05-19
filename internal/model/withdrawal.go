package model

import "time"

type Withdrawal struct {
	ID      uint
	OrderID uint
	Order   Order
	UserID  uint
	User    User
	Sum     float64
	Created time.Time
}

type WithdrawnInput struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}

type WithdrawalOutput struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	CreatedAt   time.Time `json:"processed_at"`
}
