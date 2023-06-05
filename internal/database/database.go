package database

import (
	"context"

	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type Database interface {
	CreateNewUser(ctx context.Context, user *model.User) error
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, userID uint) (*model.User, error)
	GetUserByToken(ctx context.Context, token string) (*model.User, error)
	CreateNewSession(ctx context.Context, session *model.Session) error
	GetOrderByNumber(ctx context.Context, number string) (*model.Order, error)
	CreateNewUserOrder(ctx context.Context, userID uint, number string) error
	UpdateUserBalanceAndOrder(order *model.Order, accrualInfo model.OrderAccrualInfo) error
	GetUserOrders(ctx context.Context, userID uint) ([]model.Order, error)
	GetNotReadyUserOrders(ctx context.Context, userID uint) ([]model.Order, error)
	CreateWithdraw(ctx context.Context, user *model.User, withdrawInfo model.WithdrawnInput) error
	GetUserWithdrawals(ctx context.Context, userID uint) ([]model.Withdrawal, error)
}

func GetDatabase(cfg config.DatabaseConfig) (Database, error) {
	return getPG(cfg)
}
