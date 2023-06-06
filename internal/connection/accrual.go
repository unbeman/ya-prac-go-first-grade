package connection

import (
	"context"
	"io"
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
	url := utils.FormatOrderAccrualURL(ac.address, orderNumber)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	data, err := ac.processGetOrderAccrualWithRetries(ctx, request)
	if err != nil {
		return
	}
	return utils.GetAccrualInfoFromData(data)
}

func (ac *AccrualConnection) processGetOrderAccrualWithRetries(
	ctx context.Context,
	request *http.Request,
) (data []byte, err error) {
	for try := 0; try < ac.maxRetryCount; try++ {
		ac.rateLimiter.Wait(ctx)
		response, reqErr := ac.client.Do(request)
		if reqErr != nil {
			err = reqErr
			log.Errorf("got error during request: %v", err)
			response.Body.Close()
			time.Sleep(ac.retryAfterTime)
			continue
		}

		data, err = io.ReadAll(response.Body)
		if err != nil {
			log.Error(err)
			response.Body.Close()
			return
		}
		response.Body.Close()

		isMustRetry, duration, retryErr := ac.mustRetry(data, response.StatusCode, response.Header)
		if retryErr != nil {
			err = retryErr
			log.Errorf("error occured while processing request to (%v): %v", request.URL, err)
		}
		if !isMustRetry {
			break
		}
		log.Debugf("wait (%v) for retrying request to (%v)", duration, request.URL)
		time.Sleep(duration)
	}
	return
}

func (ac *AccrualConnection) updateRequestLimit(limit int) {
	if ac.rateLimiter.Burst() != limit {
		ac.rateLimiter.SetBurst(limit)
	}
}

func (ac *AccrualConnection) mustRetry(
	data []byte,
	statusCode int,
	header http.Header,
) (isRetryNeed bool, waitDuration time.Duration, err error) {
	switch statusCode {
	case http.StatusOK:
		return false, 0, nil
	case http.StatusNoContent:
		return false, 0, apperrors.ErrNoAccrualInfo
	case http.StatusTooManyRequests:
		reqLimit, err := utils.GetRequestLimit(data)
		if err != nil {
			return false, 0, err
		}
		ac.updateRequestLimit(reqLimit)
		duration, err := utils.GetRetryWaitDuration(header)
		if err != nil {
			return false, 0, err
		}
		return true, duration, nil
	case http.StatusInternalServerError:
		return true, ac.retryAfterTime, apperrors.ErrAccrualServiceUnavailable
	default:
		log.Errorf("unexpected status code: %v", statusCode)
		return false, 0, apperrors.ErrAccrualConnection
	}
}
