package dat

import (
	"context"
	"fmt"
	"time"

	"github.com/etzba/pggo/pkg/env"
	"github.com/etzba/pggo/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

const pingRetries = 10

// Context is database context
type Context struct {
	ctx            context.Context
	Logger         *logger.Log
	ConnectionPool *pgxpool.Pool
}

// GetDatabaseContext with pgx pool and prepare for postgres connection and actions
func GetDatabaseContext() (*Context, error) {
	logger := logger.New()
	ctx := context.Background()
	cfg, err := pgxpool.ParseConfig(getConnectionString())
	if err != nil {
		logger.Error("failed to parse config", err)
	}
	cfg.MaxConns = 5
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error("failed to create pgx pool", err)
	}

	for i := 0; i < pingRetries; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			logger.Error("failed ping", err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	logger.Info("connected to database. pgx pool created")

	return &Context{
		ctx:            ctx,
		Logger:         logger,
		ConnectionPool: pool,
	}, nil
}

// getConnectionString prepare postgres connection string
func getConnectionString() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", env.PostgresUser, env.PostgresPass, env.PostgresHost, env.PostgresPort, env.PostgresDB)
}
