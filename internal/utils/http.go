package utils

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const LimitTextWordsCount = 8
const RequestLimitIndex = 3

func GetRequestLimit(response *http.Response) (reqLimit int, err error) {
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return reqLimit, err
	}
	//No more than N requests per minute allowed
	answer := strings.Split(string(data), " ") //todo reqexp?
	if len(answer) != LimitTextWordsCount {
		err = fmt.Errorf("invalid words count")
	}
	reqCount64, err := strconv.ParseInt(answer[RequestLimitIndex], 10, 0)
	if err != nil {
		return reqLimit, err
	}
	reqLimit = int(reqCount64)
	response.Body.Close()
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
