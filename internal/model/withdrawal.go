package model

import "time"

type Withdrawal struct {
	ID        uint      `json:"-"`
	Order     string    `json:"order" gorm:"unique"`
	UserID    uint      `json:"-"`
	User      User      `json:"-"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}

type WithdrawnInput struct {
	OrderNumber string  `json:"order"`
	Sum         float64 `json:"sum"`
}
