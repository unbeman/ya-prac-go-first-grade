package model

import "database/sql/driver"

//NEW — заказ загружен в систему, но не попал в обработку;
//PROCESSING — вознаграждение за заказ рассчитывается;
//INVALID — система расчёта вознаграждений отказала в расчёте;
//PROCESSED — данные по заказу проверены и информация о расчёте успешно получена.

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

func (st *OrderStatus) Scan(value interface{}) error {
	*st = OrderStatus(value.(string))
	return nil
}

func (st OrderStatus) Value() (driver.Value, error) {
	return string(st), nil
}
