package utils

import (
	"strconv"

	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"
)

func CheckOrderNumber(orderNumber string) error {
	number, err := strconv.Atoi(orderNumber)
	if err != nil {
		return err
	}
	if !valid(number) {
		return errors2.ErrInvalidOrderNumberFormat
	}
	return nil
}

func valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
