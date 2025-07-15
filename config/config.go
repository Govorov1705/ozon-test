package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
)

const (
	ModeDev           = "dev"
	ModeProd          = "prod"
	StorageInmemory   = "inmemory"
	StoragePostgreSQL = "postgresql"
)

type Config struct {
	Mode           string   `env:"MODE"`
	SecretKey      string   `env:"SECRET_KEY"`
	DBURL          string   `env:"DB_URL"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
	Storage        string   `env:"STORAGE"`
}

var Cfg Config

func InitConfig() {
	err := env.Parse(&Cfg)
	if err != nil {
		fmt.Printf("Error initializing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Config initialized")
}
