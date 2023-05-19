package model

import "time"

type Session struct {
	ID        uint
	UserID    uint
	User      User
	Token     string //todo: index
	CreatedAt time.Time
	ExpireAt  time.Time
}
