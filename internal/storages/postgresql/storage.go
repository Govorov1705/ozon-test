package postgresql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Govorov1705/ozon-test/config"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func NewStorage(DBURL string) *Storage {
	dbpool, err := pgxpool.New(context.Background(), DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	logger.Logger.Info("DB pool created")

	var absMigrationsPath string

	if config.Cfg.Mode == config.ModeDev {
		absMigrationsPath, err = filepath.Abs("./internal/storages/postgresql/migrations")
		if err != nil {
			logger.Logger.Fatal("Error converting migrations path to absolute", zap.Error(err))
		}
	} else {
		absMigrationsPath, err = filepath.Abs("./migrations")
		if err != nil {
			logger.Logger.Fatal("Error converting migrations path to absolute", zap.Error(err))
		}
	}

	migrationSourceURL := fmt.Sprintf("file://%s", absMigrationsPath)

	m, err := migrate.New(
		migrationSourceURL,
		DBURL+"sslmode=disable")
	if err != nil {
		logger.Logger.Fatal("Error creating Migrate instance", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Logger.Fatal("Error applying migrations", zap.Error(err))
	}
	logger.Logger.Info("Migrations applied")

	return &Storage{
		Pool: dbpool,
	}
}
