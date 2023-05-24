package controller

import (
	"errors"
	"fmt"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"

	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"
)

type PointsController struct {
	db *model.PG
}

func GetPointsController(db *model.PG) *PointsController {
	return &PointsController{db: db}
}

func (c PointsController) Ping() bool {
	return c.db.Ping()
}

func (c PointsController) AddUserOrder(user model.User, orderNumber string) (isNewOrder bool, err error) {
	err = utils.CheckOrderNumber(orderNumber)
	if err != nil {
		return false, errors2.ErrInvalidOrderNumberFormat
	}

	var existingOrder model.Order
	result := c.db.Conn.First(&existingOrder, "number = ?", orderNumber)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newOrder := model.Order{User: user, Status: model.StatusNew, Number: orderNumber}
		result = c.db.Conn.Create(&newOrder)
		if result.Error != nil {
			return false, fmt.Errorf("%w: %v", errors2.ErrDb, result.Error)
		}
		return true, nil
	}
	if result.Error != nil {
		return false, fmt.Errorf("%w: %v", errors2.ErrDb, result.Error)
	}
	if existingOrder.UserID != user.ID { //заказ загружен другим пользователем
		return false, errors2.ErrAlreadyExists
	}
	return false, nil
}

func (c PointsController) GetUserOrders(user model.User) (orders []model.Order, err error) {
	result := c.db.Conn.Find(&orders, "user_id = ?", user.ID).Order("created_at ASC")
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err = errors2.ErrNoRecords
		return
	}
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDb, result.Error)
	}
	log.Debugf("%+v", orders[0])
	log.Debugf("status? (%v)", orders[0].Status)
	return
}

func (c PointsController) GetUserBalance(user model.User) (balance model.UserBalanceOutput, err error) {
	balance.Withdrawn = user.Withdrawn
	balance.Current = user.CurrentBalance
	return
}

func (c PointsController) CreateWithdraw(user model.User, withdrawInfo model.WithdrawnInput) error {
	err := utils.CheckOrderNumber(withdrawInfo.OrderNumber)
	if err != nil {
		return errors2.ErrInvalidOrderNumberFormat
	}
	//todo: check order and withdraw already exists?
	err = c.db.Conn.Transaction(func(tx *gorm.DB) (txErr error) {
		result := tx.Where("id = ? and current_balance >= ?", user.ID, withdrawInfo.Sum).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors2.ErrNotEnoughPoints
		}
		user.CurrentBalance -= withdrawInfo.Sum
		user.Withdrawn += withdrawInfo.Sum
		if txErr = tx.Save(&user).Error; txErr != nil {
			return
		}
		withdraw := model.Withdrawal{Order: withdrawInfo.OrderNumber, Sum: withdrawInfo.Sum, User: user}
		if txErr = tx.Create(&withdraw).Error; txErr != nil {
			return
		}

		return
	})
	if errors.Is(err, errors2.ErrNotEnoughPoints) {
		return err
	}
	if err != nil {
		return fmt.Errorf("%w: %v", errors2.ErrDb, err)
	}
	return nil
}

func (c PointsController) GetUserWithdrawals(user model.User) (withdrawals []model.Withdrawal, err error) {
	result := c.db.Conn.Find(&withdrawals, "user_id = ?", user.ID).Order("created_at ASC")
	if result.Error != nil {
		err = fmt.Errorf("%w: %v", errors2.ErrDb, err)
		return
	}
	return
}
