package logger

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type loggerConfig struct {
	Mode string `env:"MODE"`
}

var loggerCfg loggerConfig

var Logger *zap.Logger

func InitLogger() {
	var err error

	err = env.Parse(&loggerCfg)
	if err != nil {
		log.Fatalf("Error parsing env for logger config: %v\n", err)
	}

	switch loggerCfg.Mode {
	case "dev":
		Logger, err = zap.NewDevelopment()
	case "prod":
		Logger, err = zap.NewProduction()
	default:
		Logger = zap.NewNop()
		fmt.Println("MODE env not set, using default no-op logger for testing")
	}

	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	Logger.Info("Logger initialized")
}
