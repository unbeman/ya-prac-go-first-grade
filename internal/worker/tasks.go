package worker

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type TaskOutput struct {
	Order model.Order
	Err   error
}

type TaskInterface interface {
	Do()
}

type AccrualTask struct {
	ctx    context.Context
	order  model.Order
	doFunc func(ctx context.Context, order model.Order) (model.Order, error)
	output chan TaskOutput
}

func NewAccrualTask(
	ctx context.Context,
	order model.Order,
	updateFunc func(ctx context.Context, order model.Order) (model.Order, error),
	out chan TaskOutput,
) AccrualTask {
	return AccrualTask{ctx: ctx, order: order, doFunc: updateFunc, output: out}
}

func (t AccrualTask) Do() {
	order, err := t.doFunc(t.ctx, t.order)
	if err != nil {
		log.Errorf("AccrualTask.Do got error: %v", err)
	}
	t.output <- TaskOutput{Order: order, Err: err}
	log.Infof("Task for (%v) order ended", order.Number)
}
