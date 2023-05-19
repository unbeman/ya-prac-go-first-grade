package main

import (
	"context"

	"github.com/unbeman/ya-prac-go-first-grade/internal/app"
	"github.com/unbeman/ya-prac-go-first-grade/internal/logging"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	cfg, err := app.GetConfig()
	if err != nil {
		log.Error("Can't get config:", err)
		return
	}

	logging.InitLogger(cfg.Logger)

	appl, err := app.GetApplication(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	appl.Run()
	appl.Shutdown(ctx)
}
