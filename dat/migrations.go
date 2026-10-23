package dat

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

// InitDB perform initial db migrations
func (c *Context) InitDB() error {
	tx, err := c.ConnectionPool.BeginTx(c.ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		c.Logger.Error("InitDB: failed to begin transacation", err)
		return err
	}
	// Safe to call after Commit: it returns pgx.ErrTxClosed, which we ignore.
	defer func() {
		if err = tx.Rollback(c.ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			c.Logger.Error("failed to rollback transaction", err)
		}
	}()

	tag, err := tx.Exec(c.ctx, dbMigrations)
	if err != nil {
		c.Logger.Error("InitDB: failed to execute transacation", err)
		return err
	}

	if err := tx.Commit(c.ctx); err != nil {
		c.Logger.Error("InitDB: failed to commit transacation", err)
		return err
	}

	c.Logger.Info(fmt.Sprintf("Database migration completed. Command tag: %s", tag.String()))
	return nil
}

const dbMigrations = `
DO $$ DECLARE
BEGIN
--
-- migrations pattern: if exist move on
--
IF EXISTS(SELECT 1 FROM pg_tables WHERE tablename = 'migrations') THEN
  RAISE NOTICE 'migrations table exists, skipping initial table creation';
  RETURN;
END IF;

--
-- Name: migrations; Type: TABLE; Schema: public;
--
CREATE TABLE migrations (
    name text PRIMARY KEY,
    time TIMESTAMP DEFAULT NOW()
);

END $$;

--
-- Name: locations table; Type: TABLE; Description: locations table
--
DO $$ BEGIN
IF EXISTS(SELECT 1 FROM migrations WHERE name = 'create-locations-table') THEN RETURN;
END IF;

CREATE TABLE locations (
	id SERIAL PRIMARY KEY NOT NULL,
    created timestamptz NOT NULL DEFAULT NOW(),
	name TEXT NOT NULL,
    address TEXT NOT NULL,
	longtitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL
);

INSERT INTO migrations (name) VALUES ('create-locations-table');
END $$;
`
