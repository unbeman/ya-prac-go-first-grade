package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v8"
)

const (
	ServerAddressDefault          = "127.0.0.1:8090"
	IncomingRequestTimeoutDefault = 10 * time.Second
	DatabaseURIDefault            = "postgresql://postgres:1211@localhost:5432/fgrad"
	AccrualServerAddressDefault   = "http://127.0.0.1:8080"
	LogLevelDefault               = "debug"
	TokenLifeTimeDefault          = 1 * time.Hour
	ClientTimeoutDefault          = 5 * time.Second
	RequestTimeoutDefault         = 2 * time.Second
	RateLimitDefault              = 1 * time.Minute
	RateTokensNumberDefault       = 300
	MaxRetryCountDefault          = 3
	RetryAfterTimeDefault         = 1 * time.Second
	WorkersCountDefault           = 10
	TasksSizeDefault              = 30
)

type LoggerConfig struct {
	Level string
}

type AuthConfig struct {
	TokenLifeTime time.Duration
}

type AccrualConnConfig struct {
	ServerAddress   string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	ClientTimeout   time.Duration
	RequestTimeout  time.Duration
	RateLimit       time.Duration
	RateTokensCount int
	MaxRetryCount   int
	RetryAfterTime  time.Duration
}

type ApplicationConfig struct {
	ServerAddress          string `env:"RUN_ADDRESS"`
	IncomingRequestTimeout time.Duration
	Database               DatabaseConfig
	Logger                 LoggerConfig
	Auth                   AuthConfig
	AccrualConn            AccrualConnConfig
	WorkerPool             WorkerPoolConfig
}

type WorkerPoolConfig struct {
	WorkersCount int
	TasksSize    int
}

type DatabaseConfig struct {
	DatabaseURI string `env:"DATABASE_URI"`
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
	cfg.Database.DatabaseURI = *dbURI
	cfg.AccrualConn.ServerAddress = *accrualAddr
	cfg.Logger.Level = *logLevel
	return nil
}

func GetConfig() (ApplicationConfig, error) {
	cfg := ApplicationConfig{
		ServerAddress:          ServerAddressDefault,
		IncomingRequestTimeout: IncomingRequestTimeoutDefault,

		Logger: LoggerConfig{
			Level: LogLevelDefault,
		},
		Auth: AuthConfig{
			TokenLifeTime: TokenLifeTimeDefault,
		},
		AccrualConn: AccrualConnConfig{
			ServerAddress:   AccrualServerAddressDefault,
			ClientTimeout:   ClientTimeoutDefault,
			RequestTimeout:  RequestTimeoutDefault,
			RateLimit:       RateLimitDefault,
			RateTokensCount: RateTokensNumberDefault,
			MaxRetryCount:   MaxRetryCountDefault,
			RetryAfterTime:  RetryAfterTimeDefault,
		},
		WorkerPool: WorkerPoolConfig{
			WorkersCount: WorkersCountDefault,
			TasksSize:    TasksSizeDefault,
		},
		Database: DatabaseConfig{
			DatabaseURI: DatabaseURIDefault,
		},
	}
	if err := cfg.parseFlags(); err != nil {
		return cfg, err
	}
	if err := cfg.parseEnv(); err != nil {
		return cfg, err
	}
	return cfg, nil
}
