package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/unbeman/ya-prac-go-first-grade/internal/controller"
	"github.com/unbeman/ya-prac-go-first-grade/internal/handler"
	"github.com/unbeman/ya-prac-go-first-grade/internal/model"

	log "github.com/sirupsen/logrus"
)

type application struct {
	server http.Server
}

func GetApplication(cfg ApplicationConfig) (*application, error) { //TODO: sync once
	db, err := model.GetDatabase(cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	authControl := controller.GetAuthController(db)
	pointsControl := controller.GetPointsController(db)
	hndlr := handler.GetAppHandler(authControl, pointsControl)
	return &application{
		server: http.Server{Addr: cfg.ServerAddress, Handler: hndlr},
	}, nil
}

func (a *application) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := a.server.ListenAndServe()
		log.Error(err)
		wg.Done()
	}()
	log.Info("Http server started")
	wg.Wait()
}

func (a *application) Shutdown(ctx context.Context) {
	a.server.Shutdown(ctx)
	log.Info("Http server closed")
}
