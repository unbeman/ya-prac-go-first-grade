package logging

import (
	_ "github.com/golang/mock/mockgen/model"
	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-go-first-grade/internal/config"
)

const (
	LogDebug = "debug"
	LogInfo  = "info"
)

func InitLogger(cfg config.LoggerConfig) {
	switch cfg.Level {
	case LogInfo:
		log.SetLevel(log.InfoLevel)
	case LogDebug:
		log.SetLevel(log.DebugLevel)
	}
}
