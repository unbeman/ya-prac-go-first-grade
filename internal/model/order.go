package model

import "time"

type Order struct {
	ID        uint        `json:"-"`
	UserID    uint        `json:"-"`
	User      User        `json:"-"`
	Status    OrderStatus `json:"status" gorm:"type:order_status;default:'NEW'"`
	Number    string      `json:"order" gorm:"unique"`
	Accrual   float64     `json:"accrual,omitempty"`
	CreatedAt time.Time   `json:"uploaded_at"`
}

type OrderAccrualInfo struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
