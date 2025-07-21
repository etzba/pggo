package dat

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// InitDB perform initial db migrations
func (c *Context) InitDB() error {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		c.Logger.Error("could not open connection to database", err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			c.Logger.Error("failed to close db connection", err)
		}
	}()

	if err = db.Ping(); err != nil {
		c.Logger.Error("could not ping database", err)
		return err
	}

	_, err = db.Exec(dbMigrations)
	if err != nil {
		c.Logger.Error("could not run db migrations", err)
		return err
	}

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
