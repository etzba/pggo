package dat

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func (c *Context) InitDB() error {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		c.Logger.Error("could not open connection to database", err)
		return err
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		c.Logger.Error("could not ping database", err)
		return err
	}

	_, err = db.Exec(dbMigrations)
	if err != nil {
		c.Logger.Error("could not run db migrations", err)
		return err
	}
	c.Logger.Info("db migration completed")
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
`
