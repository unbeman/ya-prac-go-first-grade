package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"

	errors2 "github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type AccrualConnection struct {
	client         http.Client
	address        string
	requestTimeout time.Duration
	rateLimiter    *rate.Limiter
}

func NewAccrualConnection(cfg config.AccrualConnConfig) *AccrualConnection {
	cli := http.Client{Timeout: cfg.ClientTimeout}
	rl := rate.NewLimiter(rate.Every(cfg.RateLimit), cfg.RateTokensNumber) //не больше RateTokensNumber запросов в RateLimit
	return &AccrualConnection{
		client:         cli,
		address:        cfg.ServerAddress,
		requestTimeout: cfg.RequestTimeout,
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
	log.Debugf("GetOrderAccrual: recieved status %v, request GET: %v", response.StatusCode, url)
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
