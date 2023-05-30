package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"

	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/app-errors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

//todo: понятное название пакета

// todo: move to config
const (
	ClientTimeoutDefault  = 5 * time.Second
	RequestTimeoutDefault = 2 * time.Second
	RateLimitDefault      = 60
)

type AccrualConnection struct {
	client         http.Client
	address        string
	requestTimeout time.Duration
	rateLimiter    *rate.Limiter
}

func NewAccrualConnection(addr string) *AccrualConnection {
	cli := http.Client{Timeout: ClientTimeoutDefault}
	rl := rate.NewLimiter(rate.Every(1*time.Minute), RateLimitDefault) //не больше rateLimit запросов в минуту
	return &AccrualConnection{
		client:         cli,
		address:        addr,
		requestTimeout: RequestTimeoutDefault,
		rateLimiter:    rl,
	}
}

// todo: wrap errors
func (ac *AccrualConnection) GetOrderAccrual(ctx context.Context, orderNumber string) (orderInfo model.OrderAccrualInfo, err error) {
	ac.rateLimiter.Wait(ctx)
	ctx2, cancel := context.WithTimeout(ctx, ac.requestTimeout)
	defer cancel()
	url := fmt.Sprintf("%v/api/orders/%v", ac.address, orderNumber)
	request, err := http.NewRequestWithContext(ctx2, http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := ac.client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("GetOrderAccrual: recieved status %v, request GET: %v", response.StatusCode, url)
	switch response.StatusCode {
	case http.StatusOK:
		defer response.Body.Close()
		jD := json.NewDecoder(response.Body)
		err = jD.Decode(&orderInfo)
		if err != nil {
			log.Error(err)
			return
		}
		return orderInfo, nil
	case http.StatusNoContent:
		err = errors2.ErrNoAccrualInfo
	case http.StatusTooManyRequests: //не воспроизводится
		request.Header.Get("Retry-After")
		//todo: go retry
	case http.StatusInternalServerError:
		err = errors2.ErrNoAccrualInfo //todo another err
	}
	return
}
