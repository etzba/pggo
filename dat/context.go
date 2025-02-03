package dat

import (
	"context"
	"fmt"

	"github.com/etzba/pggo/pkg/env"
	"github.com/etzba/pggo/pkg/logger"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Context struct {
	ctx            context.Context
	Logger         *logger.Log
	ConnectionPool *pgxpool.Pool
}

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

func (c *Context) aquireConn() {
	conn, err := c.ConnectionPool.Acquire(context.Background())
	if err != nil {
		c.Logger.Error("could not aquire connection to database", err)
	}
	defer conn.Release()
}

func getConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", env.PostgresUser, env.PostgresPass, env.PostgresHost, env.PostgresPort, env.PostgresDB)
}
