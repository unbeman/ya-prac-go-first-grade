package controller

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/connection"
	"github.com/unbeman/ya-prac-go-first-grade/internal/database"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
	"github.com/unbeman/ya-prac-go-first-grade/internal/worker"
)

type PointsController struct {
	db                database.Database
	accrualConnection connection.AccrualConnector
	wp                *worker.WorkersPool
}

func GetPointsController(db database.Database, accConn connection.AccrualConnector, wp *worker.WorkersPool) *PointsController {
	return &PointsController{db: db, accrualConnection: accConn, wp: wp}
}

func (c PointsController) AddUserOrder(ctx context.Context, user *model.User, orderNumber string) (isNewOrder bool, err error) {
	err = utils.CheckOrderNumber(orderNumber)
	if err != nil {
		return false, apperrors.ErrInvalidOrderNumberFormat
	}

	existingOrder, err := c.db.GetOrderByNumber(ctx, orderNumber)
	if errors.Is(err, apperrors.ErrNoRecords) {
		err = c.db.CreateNewUserOrder(ctx, user.ID, orderNumber)
		if err != nil {
			return
		}
		return true, nil
	}
	if err != nil {
		return
	}
	if existingOrder.UserID != user.ID { //заказ загружен другим пользователем
		return false, apperrors.ErrAlreadyExists
	}
	return false, nil
}

func (c PointsController) updateUserOrder(order model.Order) (model.Order, error) {
	orderAccrualInfo, err := c.accrualConnection.GetOrderAccrual(context.TODO(), order.Number)
	if err != nil {
		return order, err
	}
	log.Debug("updateUserOrder orderAccrualInfo: ", orderAccrualInfo)
	err = c.db.UpdateUserBalanceAndOrder(&order, orderAccrualInfo)
	return order, err
}

func (c PointsController) UpdateUserOrders(ctx context.Context, user *model.User) error {
	notReadyOrders, err := c.db.GetNotReadyOrders(ctx, user.ID)
	if err != nil {
		log.Error(err)
		return err
	}

	taskOutput := make(chan worker.TaskOutput, len(notReadyOrders))
	for _, order := range notReadyOrders {
		updateOrder := &worker.Task{Order: order, DoFunc: c.updateUserOrder, OutputErr: taskOutput}
		c.wp.AddTask(updateOrder)
	}
	for idx := 0; idx < len(notReadyOrders); idx++ {
		if out := <-taskOutput; out.Err != nil {
			log.Error(out.Err)
		}
	}
	return nil
}

func (c PointsController) GetUserOrders(ctx context.Context, user *model.User) (orders []model.Order, err error) {
	return c.db.GetUserOrders(ctx, user.ID)
}

func (c PointsController) GetUserBalance(ctx context.Context, user *model.User) (balance model.UserBalanceOutput, err error) {
	user, err = c.db.GetUserByID(ctx, user.ID)
	if err != nil {
		return
	}
	balance.Withdrawn = user.Withdrawn
	balance.Current = user.CurrentBalance
	return
}

func (c PointsController) CreateWithdraw(ctx context.Context, user *model.User, withdrawInfo model.WithdrawnInput) error {
	err := utils.CheckOrderNumber(withdrawInfo.OrderNumber)
	if err != nil {
		return apperrors.ErrInvalidOrderNumberFormat
	}
	err = c.db.CreateWithdraw(ctx, user, withdrawInfo)
	return err
}

func (c PointsController) GetUserWithdrawals(ctx context.Context, user *model.User) (withdrawals []model.Withdrawal, err error) {
	return c.db.GetUserWithdrawals(ctx, user.ID)
}
