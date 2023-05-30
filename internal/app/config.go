package app

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

const (
	ServerAddressDefault        = "127.0.0.1:8090"
	DatabaseURIDefault          = "postgresql://postgres:1211@localhost:5432/fgrad"
	AccrualServerAddressDefault = "127.0.0.1:8080"
	LogLevelDefault             = "info"
)

type LoggerConfig struct {
	Level string
}

type ApplicationConfig struct {
	ServerAddress        string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualServerAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Logger               LoggerConfig
}

func (cfg *ApplicationConfig) parseEnv() error {
	if err := env.Parse(cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *ApplicationConfig) parseFlags() error {
	serverAddr := flag.String("a", ServerAddressDefault, "server address")
	dbURI := flag.String("d", DatabaseURIDefault, "database address")
	accrualAddr := flag.String("r", AccrualServerAddressDefault, "accrual server address")
	logLevel := flag.String("l", LogLevelDefault, "log level, allowed: info, debug")
	flag.Parse()
	cfg.ServerAddress = *serverAddr
	cfg.DatabaseURI = *dbURI
	cfg.AccrualServerAddress = *accrualAddr
	cfg.Logger.Level = *logLevel
	return nil
}

func GetConfig() (ApplicationConfig, error) {
	cfg := ApplicationConfig{
		ServerAddress:        ServerAddressDefault,
		DatabaseURI:          DatabaseURIDefault,
		AccrualServerAddress: AccrualServerAddressDefault,
		Logger:               LoggerConfig{Level: LogLevelDefault},
	}
	if err := cfg.parseFlags(); err != nil {
		return cfg, err
	}
	if err := cfg.parseEnv(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
