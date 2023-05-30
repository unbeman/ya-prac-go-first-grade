package app

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/accrual"
	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/handler"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"
)

type application struct {
	server        http.Server
	pointsControl *controller.PointsController
}

func GetApplication(cfg ApplicationConfig) (*application, error) { //TODO: sync once
	db, err := model.GetDatabase(cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	accConn := accrual.NewAccrualConnection(cfg.AccrualServerAddress)
	authControl := controller.GetAuthController(db)
	orderChan := make(chan string)
	pointsControl := controller.GetPointsController(db, orderChan, accConn)
	hndlr := handler.GetAppHandler(authControl, pointsControl)
	return &application{
		server:        http.Server{Addr: cfg.ServerAddress, Handler: hndlr},
		pointsControl: pointsControl,
	}, nil
}

func (a *application) Run() {
	var wg sync.WaitGroup

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
}
