package logging

import (
	"github.com/unbeman/ya-prac-go-first-grade/internal/app"

	log "github.com/sirupsen/logrus"
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
