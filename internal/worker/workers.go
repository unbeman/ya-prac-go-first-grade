package worker

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

//todo intefraces

type TaskOutput struct {
	Err error
}

type Task struct {
	Order     model.Order
	DoFunc    func(model.Order) (model.Order, error)
	OutputErr chan TaskOutput
}

func (t *Task) Do() {
	_, err := t.DoFunc(t.Order)
	if err != nil {
		log.Error(err)
	}
	t.OutputErr <- TaskOutput{Err: err}
}

type WorkersPool struct {
	wokersCount int
	tasks       chan *Task
	tasksSize   int
	waitGroup   sync.WaitGroup
}

func NewWorkersPool(wokersCount int) *WorkersPool {
	return &WorkersPool{
		wokersCount: wokersCount,
		tasks:       make(chan *Task, 10), //todo buffer
	}
}

func (wp *WorkersPool) Run() {
	log.Infof("starting worker pool %d workers", wp.wokersCount)
	for idx := 0; idx < wp.wokersCount; idx++ {
		wp.waitGroup.Add(1)
		go func(idx int, tasks chan *Task) {
			log.Infof("worker %d started", idx)
			defer wp.waitGroup.Done()
			for task := range tasks {
				task.Do()
			}
			log.Infof("worker %d: finished", idx)
		}(idx, wp.tasks)
	}
	wp.waitGroup.Wait()
}

func (wp *WorkersPool) AddTask(task *Task) {
	wp.tasks <- task
}

func (wp *WorkersPool) Shutdown() {
	close(wp.tasks)
}
