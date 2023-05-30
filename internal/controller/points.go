package controller

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/accrual"
	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

type PointsController struct {
	db                *model.PG
	accrualConnection *accrual.AccrualConnection
}

func GetPointsController(db *model.PG, accConn *accrual.AccrualConnection) *PointsController {
	return &PointsController{db: db, accrualConnection: accConn}
}

func (c PointsController) Ping() bool {
	return c.db.Ping()
}

func (c PointsController) AddUserOrder(user *model.User, orderNumber string) (isNewOrder bool, err error) {
	err = utils.CheckOrderNumber(orderNumber)
	if err != nil {
		return false, errors2.ErrInvalidOrderNumberFormat
	}

	existingOrder, err := c.db.GetOrderByNumber(orderNumber)
	if errors.Is(err, errors2.ErrNoRecords) {
		err = c.db.CreateNewUserOrder(user.ID, orderNumber)
		if err != nil {
			return
		}
		return true, nil
	}
	if err != nil {
		return
	}
	if existingOrder.UserID != user.ID { //заказ загружен другим пользователем
		return false, errors2.ErrAlreadyExists
	}
	return false, nil
}

func (c PointsController) updateUserOrder(order *model.Order) error {
	orderAccrualInfo, err := c.accrualConnection.GetOrderAccrual(context.TODO(), order.Number)
	if err != nil {
		return err
	}
	log.Debug("updateUserOrder orderAccrualInfo: ", orderAccrualInfo)
	if order.Status != orderAccrualInfo.Status {
		order.Status = orderAccrualInfo.Status
		order.Accrual = orderAccrualInfo.Accrual
		err = c.db.UpdateUserBalanceAndOrder(order)
	}
	return err
}

func (c PointsController) GetUserOrders(user *model.User) (orders []model.Order, err error) {
	orders, err = c.db.GetUserOrders(user.ID)
	if err != nil {
		log.Error(err)
		return
	}
	for _, order := range orders {
		if order.Status == model.StatusProcessed || order.Status == model.StatusInvalid {
			continue
		}
		if updErr := c.updateUserOrder(&order); updErr != nil { //todo: async
			log.Error("GetUserOrders: ", updErr)
			if errors.Is(updErr, errors2.ErrDB) {
				err = updErr
				return
			} //ignore accrual conn error
		}
	}
	return
}

func (c PointsController) GetUserBalance(user *model.User) (balance model.UserBalanceOutput, err error) {
	//todo: проверить заказы
	balance.Withdrawn = user.Withdrawn
	balance.Current = user.CurrentBalance
	return
}

func (c PointsController) CreateWithdraw(user *model.User, withdrawInfo model.WithdrawnInput) error {
	err := utils.CheckOrderNumber(withdrawInfo.OrderNumber)
	if err != nil {
		return errors2.ErrInvalidOrderNumberFormat
	}
	//todo: проверить заказы
	err = c.db.CreateWithdraw(user, withdrawInfo)
	return err
}

func (c PointsController) GetUserWithdrawals(user *model.User) (withdrawals []model.Withdrawal, err error) {
	return c.db.GetUserWithdrawals(user.ID)
}
