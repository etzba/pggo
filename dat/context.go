package dat

import (
	"context"
	"fmt"

	"github.com/etzba/pggo/pkg/env"
	"github.com/etzba/pggo/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

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
	pool, err := pgxpool.Connect(ctx, getConnectionString())
	if err != nil {
		logger.Error("could not create a db connection pool", err)
		return nil, err
	}

	return &Context{
		ctx:            ctx,
		Logger:         logger,
		ConnectionPool: pool,
	}, nil
}

// getConnectionString prepare postgres connection string
func getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", env.PostgresUser, env.PostgresPass, env.PostgresHost, env.PostgresPort, env.PostgresDB)
}
