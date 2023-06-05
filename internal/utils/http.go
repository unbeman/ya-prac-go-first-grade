package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

const LimitTextWordsCount = 8
const RequestLimitIndex = 3

func GetRequestLimit(data []byte) (reqLimit int, err error) {
	//No more than N requests per minute allowed
	answer := strings.Split(string(data), " ") //todo reqexp?
	if len(answer) != LimitTextWordsCount {
		err = errors.New("invalid words count")
		return
	}
	reqCount64, err := strconv.ParseInt(answer[RequestLimitIndex], 10, 0)
	if err != nil {
		return reqLimit, err
	}
	reqLimit = int(reqCount64)
	return reqLimit, nil
}

func GetRetryWaitDuration(header http.Header) (duration time.Duration, err error) {
	retryAfter := header.Get("Retry-After")
	retryTimeSec, err := strconv.ParseInt(retryAfter, 10, 64)
	if err != nil {
		log.Error(err)
		return
	}
	duration = time.Duration(retryTimeSec)
	return
}

func FormatOrderAccrualURL(address, orderNumber string) string {
	return fmt.Sprintf("%v/api/orders/%v", address, orderNumber)
}

func GetAccrualInfoFromData(data []byte) (model.OrderAccrualInfo, error) {
	info := model.OrderAccrualInfo{}
	err := json.Unmarshal(data, &info)
	return info, err
}
