package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/HadeedTariq/go-production-grade-api/internal/utils/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()

	if err != nil {
		logger.Error("error loading env variables", "error", err.Error())
	}

	cfg := config{
		addr: ":3000",
		dbConfig: dbConfig{
			dsn: env.GetEnvString("GOOSE_DBSTRING", "host=localhost user=myuser password=mypassword dbname=daily-dev-db sslmode=disable"),
		},
	}

	pool, err := pgxpool.New(ctx, cfg.dbConfig.dsn)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	logger.Info("connected to database", "dsn", cfg.dbConfig.dsn)

	app := application{
		config: cfg,
		db:     pool,
	}

	handler := app.mount()

	err = app.run(handler)

	if err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
