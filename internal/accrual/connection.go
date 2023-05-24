package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

//todo: понятное название пакета

// todo: move to config
const (
	ClientTimeoutDefault  = 5 * time.Second
	RequestTimeoutDefault = 2 * time.Second
	RateLimitDefault      = 10
)

type accrualConnection struct {
	client         http.Client
	address        string
	requestTimeout time.Duration
	rateLimiter    *rate.Limiter
}

func NewAccrualConnection(addr string) *accrualConnection {
	cli := http.Client{Timeout: ClientTimeoutDefault}
	rl := rate.NewLimiter(rate.Every(1*time.Minute), RateLimitDefault) //не больше rateLimit запросов в минуту
	return &accrualConnection{
		client:         cli,
		address:        addr,
		requestTimeout: RequestTimeoutDefault,
		rateLimiter:    rl,
	}
}

// todo: wrap errors
func (ac *accrualConnection) GetOrderAccrual(ctx context.Context, orderNumber string) (orderInfo model.OrderAccrualInfo, err error) {
	ctx2, cancel := context.WithTimeout(ctx, ac.requestTimeout)
	defer cancel()
	url := fmt.Sprintf("%v/api/orders/%v", ac.address, orderNumber)
	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, nil)
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
	switch response.StatusCode {
	case http.StatusOK:
	case http.StatusNoContent:
	case http.StatusTooManyRequests:
		request.Header.Get("Retry-After")
	case http.StatusInternalServerError:
	}
	defer response.Body.Close()
	jD := json.NewDecoder(response.Body)
	err = jD.Decode(&orderInfo)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
