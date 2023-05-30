package model

import (
	"fmt"
	"time"
)

type User struct {
	ID             uint   `gorm:"primary_key"` //todo: use uuid
	Login          string `gorm:"uniqueIndex"`
	HashPassword   string
	CurrentBalance float64
	Withdrawn      float64
	CreatedAt      time.Time
}

func (u User) String() string {
	return fmt.Sprintf("ID=%v;Login=%v;HashPassword=%v;CurrentBalance=%v;Withdrawn=%v;CreatedAt=%v;",
		u.ID, u.Login, u.HashPassword, u.CurrentBalance, u.Withdrawn, u.CreatedAt)
}

type UserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserBalanceOutput struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
