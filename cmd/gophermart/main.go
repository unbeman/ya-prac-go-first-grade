package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/app"
	"github.com/unbeman/ya-prac-go-first-grade/internal/logging"
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
	log.Info(cfg)
	appl, err := app.GetApplication(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(
			exit,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)

		for {
			sig := <-exit
			log.Infof("Got signal '%v'", sig)
			appl.Shutdown(ctx)
		}
	}()

	appl.Run()
	//appl.Shutdown(ctx)
}
