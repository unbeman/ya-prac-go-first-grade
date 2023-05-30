package logging

import (
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-first-grade/internal/app"
)

const (
	LogDebug = "debug"
	LogInfo  = "info"
)

func InitLogger(cfg app.LoggerConfig) {
	switch cfg.Level {
	case LogInfo:
		log.SetLevel(log.InfoLevel)
	case LogDebug:
		log.SetLevel(log.DebugLevel)
	}
}
