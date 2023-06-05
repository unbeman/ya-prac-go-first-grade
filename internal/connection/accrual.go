package connection

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/apperrors"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/utils"

	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type AccrualConnector interface {
	GetOrderAccrual(ctx context.Context, orderNumber string) (model.OrderAccrualInfo, error)
}

func GetAccrualConnector(cfg config.AccrualConnConfig) AccrualConnector {
	return NewAccrualConnection(cfg)
}

type AccrualConnection struct {
	client         http.Client
	address        string
	requestTimeout time.Duration
	rateLimiter    *rate.Limiter
	maxRetryCount  int
	retryAfterTime time.Duration
}

func NewAccrualConnection(cfg config.AccrualConnConfig) *AccrualConnection {
	cli := http.Client{Timeout: cfg.ClientTimeout}
	rl := rate.NewLimiter(rate.Every(cfg.RateLimit), cfg.RateTokensCount) //не больше RateTokensCount запросов в RateLimit
	return &AccrualConnection{
		client:         cli,
		address:        cfg.ServerAddress,
		requestTimeout: cfg.RequestTimeout,
		rateLimiter:    rl,
		maxRetryCount:  cfg.MaxRetryCount,
		retryAfterTime: cfg.RetryAfterTime,
	}
}

func (ac *AccrualConnection) GetOrderAccrual(ctx context.Context, orderNumber string) (orderInfo model.OrderAccrualInfo, err error) {
	ac.rateLimiter.Wait(ctx)
	url := utils.FormatOrderAccrualURL(ac.address, orderNumber)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := ac.processGetOrderAccrualWithRetries(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	log.Debugf("got HTTP status from (%v): %v", request.URL, response.StatusCode)
	orderInfo, err = utils.GetAccrualInfoFromResponse(response)
	if err != nil {
		return
	}
	return
}

func (ac *AccrualConnection) processGetOrderAccrualWithRetries(request *http.Request) (response *http.Response, err error) {
	for try := 0; try < ac.maxRetryCount; try++ {
		log.Debugf("making GET request to (%v)", request.URL)
		response, err = ac.client.Do(request)
		isMustRetry, duration, err := ac.mustRetry(response, err)
		if err != nil {
			log.Errorf("error occured while processing request to (%v): %v", request.URL, err)
		}
		if !isMustRetry {
			break
		}
		log.Debugf("wait (%v) for retrying request to (%v)", duration, request.URL)
		time.Sleep(duration)
	}
	return response, err
}

func (ac *AccrualConnection) updateRequestLimit(limit int) {
	if ac.rateLimiter.Burst() != limit {
		ac.rateLimiter.SetBurst(limit)
	}
}

func (ac *AccrualConnection) mustRetry(
	response *http.Response,
	requestErr error,
) (isRetryNeed bool, waitDuration time.Duration, err error) {
	if requestErr != nil {
		log.Errorf("got error during request: %v", requestErr)
		return true, ac.retryAfterTime, requestErr
	}
	switch response.StatusCode {
	case http.StatusOK:
		return false, 0, nil
	case http.StatusNoContent:
		return false, 0, apperrors.ErrNoAccrualInfo
	case http.StatusTooManyRequests:
		reqLimit, err := utils.GetRequestLimit(response)
		if err != nil {
			return false, 0, err
		}
		ac.updateRequestLimit(reqLimit)
		duration, err := utils.GetRetryWaitDuration(response.Header)
		if err != nil {
			return false, 0, err
		}
		return true, duration, nil
	case http.StatusInternalServerError:
		return true, ac.retryAfterTime, apperrors.ErrAccrualServiceUnavailable
	default:
		log.Errorf("unexpected status code: %v", response.StatusCode)
		return false, 0, apperrors.ErrAccrualConnection
	}
}
