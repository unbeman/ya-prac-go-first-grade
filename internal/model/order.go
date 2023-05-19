package model

import "time"

type Order struct {
	ID        uint
	UserID    uint
	User      User
	Status    OrderStatus `gorm:"type:order_status;default:'NEW'"`
	Number    string      `gorm:"unique"`
	Accrual   float64
	CreatedAt time.Time
}

type OrderInput struct {
	Number string
}

type OrderOutput struct {
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"uploaded_at"`
}
