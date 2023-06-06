package worker

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
)

type WorkersPool struct {
	wokersCount int
	tasks       chan TaskInterface
	tasksSize   int
	waitGroup   sync.WaitGroup
}

func NewWorkersPool(cfg config.WorkerPoolConfig) *WorkersPool {
	return &WorkersPool{
		wokersCount: cfg.WorkersCount,
		tasks:       make(chan TaskInterface, cfg.TasksSize),
	}
}

func (wp *WorkersPool) Run() {
	log.Infof("starting worker pool %d workers", wp.wokersCount)
	for idx := 0; idx < wp.wokersCount; idx++ {
		wp.waitGroup.Add(1)
		go func(idx int, tasks chan TaskInterface) {
			defer wp.waitGroup.Done()
			log.Infof("worker %d started", idx)
			for task := range tasks {
				log.Infof("worker %d: starting task", idx)
				task.Do()
			}
			log.Infof("worker %d: finished", idx)
		}(idx, wp.tasks)
	}
	wp.waitGroup.Wait()
}

func (wp *WorkersPool) AddTask(task TaskInterface) {
	wp.tasks <- task
}

func (wp *WorkersPool) Shutdown() {
	close(wp.tasks)
}
