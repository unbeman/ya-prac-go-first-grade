package app

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/connection"
	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/database"
	"github.com/unbeman/ya-prac-go-first-grade/internal/handler"
	"github.com/unbeman/ya-prac-go-first-grade/internal/worker"
)

type application struct {
	server        http.Server
	pointsControl *controller.PointsController
	workerPool    *worker.WorkersPool
}

func GetApplication(cfg config.ApplicationConfig) (*application, error) { //TODO: sync once
	db, err := database.GetDatabase(cfg.Database)
	if err != nil {
		return nil, err
	}
	accConn := connection.GetAccrualConnector(cfg.AccrualConn)
	authControl := controller.GetAuthController(db, cfg.Auth)
	workerPool := worker.NewWorkersPool(cfg.WorkerPool)
	pointsControl := controller.GetPointsController(db, accConn, workerPool)
	hndlr := handler.GetAppHandler(cfg.IncomingRequestTimeout, authControl, pointsControl)
	return &application{
		server:        http.Server{Addr: cfg.ServerAddress, Handler: hndlr},
		pointsControl: pointsControl,
		workerPool:    workerPool,
	}, nil
}

func (a *application) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		a.workerPool.Run()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		err := a.server.ListenAndServe()
		log.Info(err)
		wg.Done()
	}()
	log.Info("Http server started")
	wg.Wait()
}

func (a *application) Shutdown(ctx context.Context) {
	a.server.Shutdown(ctx)
	a.workerPool.Shutdown()
}
