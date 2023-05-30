package model

import (
	"fmt"
	"time"
)

type Order struct {
	ID        uint        `json:"-"`
	UserID    uint        `json:"-"`
	User      User        `json:"-"`
	Status    OrderStatus `json:"status" gorm:"type:order_status;default:'NEW'"`
	Number    string      `json:"number" gorm:"unique"`
	Accrual   float64     `json:"accrual,omitempty"`
	CreatedAt time.Time   `json:"uploaded_at"`
}

func (o Order) String() string {
	return fmt.Sprintf("ID=%v;UserID=%v;User=[%v];Status=%v;Number=%v;Accrual=%v;CreatedAt=%v",
		o.ID, o.UserID, o.User, o.Status, o.Number, o.Accrual, o.CreatedAt)
}

type OrderAccrualInfo struct {
	Number  string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual,omitempty"`
}
