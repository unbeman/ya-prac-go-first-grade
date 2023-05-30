package app

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
	"github.com/unbeman/ya-prac-go-first-grade/internal/database"

	"github.com/unbeman/ya-prac-go-first-grade/internal/accrual"
	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/handler"
)

type application struct {
	server        http.Server
	pointsControl *controller.PointsController
}

func GetApplication(cfg config.ApplicationConfig) (*application, error) { //TODO: sync once
	db, err := database.GetDatabase(cfg.DatabaseURI) //todo: interface
	if err != nil {
		return nil, err
	}
	accConn := accrual.NewAccrualConnection(cfg.AccrualConn)
	authControl := controller.GetAuthController(db, cfg.Auth)
	pointsControl := controller.GetPointsController(db, accConn)
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
