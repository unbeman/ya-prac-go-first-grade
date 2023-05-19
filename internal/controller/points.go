package controller

import (
	"errors"
	"fmt"
	"gorm.io/gorm"

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
