package dat

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

// Location
type Location struct {
	Id         int
	Created    time.Time
	Name       string
	Address    string
	Longtitude float64
	Latitude   float64
}

// Locations
type Locations interface {
	InsertLocationIntoDatabase(loc Location) error
	GetLocationDetailsByID(locationId int) (*Location, error)
	GetAllLocationDetails() (*[]Location, error)
	UpdateLocationDetails(int, Location) error
	DeleteLocationFromDatabase(int) error
}

// InsertLocationIntoDatabase
func (c *Context) InsertLocationIntoDatabase(loc Location) error {
	query := "INSERT INTO locations (name, address, longtitude, latitude) VALUES ($1, $2, $3, $4)"

	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		c.Logger.Error("Failed to begin transaction", err)
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(context.TODO()); err != nil {
				c.Logger.Error("Failed to rollback transaction", err)
			}
		} else {
			if err := tx.Commit(context.TODO()); err != nil {
				c.Logger.Error("Failed to commit transaction", err)
			}
		}
	}()

	row, err := tx.Exec(context.Background(), query, loc.Name, loc.Address, loc.Longtitude, loc.Latitude)
	if err != nil {
		c.Logger.Error("Failed to execute transaction", err)
		return err
	}

	c.Logger.Info(fmt.Sprintf("affected rows %d", row.RowsAffected()))
	return nil
}

func (c *Context) GetLocationDetailsByID(locationId int) (*Location, error) {
	query := fmt.Sprintf("SELECT * FROM locations WHERE id=%d", locationId)
	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(context.TODO()); err != nil {
				c.Logger.Error("Failed to rollback transaction", err)
			}
		} else {
			if err := tx.Commit(context.TODO()); err != nil {
				c.Logger.Error("Failed to commit transaction", err)
			}
		}
	}()

	row := tx.QueryRow(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	l := &Location{}
	if err := row.Scan(&l.Id, &l.Created, &l.Name, &l.Address, &l.Longtitude, &l.Latitude); err != nil {
		return nil, err
	}
	return l, nil
}

func (c *Context) GetAllLocationDetails() ([]Location, error) {
	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(context.TODO()); err != nil {
				c.Logger.Error("Failed to rollback transaction", err)
			}
		} else {
			if err := tx.Commit(context.TODO()); err != nil {
				c.Logger.Error("Failed to commit transaction", err)
			}
		}
	}()

	rows, err := tx.Query(context.TODO(), "SELECT * FROM locations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := []Location{}
	for rows.Next() {
		l := Location{}
		if err := rows.Scan(&l.Id, &l.Created, &l.Name, &l.Address, &l.Longtitude, &l.Latitude); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (c *Context) UpdateLocationDetails(id int, loc Location) error {
	locMap := make(map[string]interface{})
	if loc.Name != "" {
		locMap["name"] = loc.Name
	}
	if loc.Address != "" {
		locMap["address"] = loc.Address
	}
	if loc.Longtitude != 0 {
		locMap["longtitude"] = loc.Longtitude
	}
	if loc.Latitude != 0 {
		locMap["latitude"] = loc.Latitude
	}

	query := "UPDATE locations SET"
	count := 1
	for k := range locMap {
		query += fmt.Sprintf(" %s = $%d,", k, count)
		count++
	}
	query = strings.TrimSuffix(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", count)
	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		c.Logger.Error("Failed to begin transaction", err)
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(context.TODO()); err != nil {
				c.Logger.Error("Failed to rollback transaction", err)
			}
		} else {
			if err := tx.Commit(context.TODO()); err != nil {
				c.Logger.Error("Failed to commit transaction", err)
			}
		}
	}()

	rows, err := tx.Exec(context.Background(), query, loc.Name, loc.Address, loc.Longtitude, loc.Latitude, id)
	if err != nil {
		c.Logger.Error("Failed to execute transaction", err)
		return err
	}

	c.Logger.Info(fmt.Sprintf("affected rows %d", rows.RowsAffected()))
	return nil
}

func (c *Context) DeleteLocationFromDatabase(locationId int) error {
	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(context.TODO()); err != nil {
				c.Logger.Error("Failed to rollback transaction", err)
			}
		} else {
			if err := tx.Commit(context.TODO()); err != nil {
				c.Logger.Error("Failed to commit transaction", err)
			}
		}
	}()

	_, err = tx.Query(context.Background(), "DELETE FROM locations WHERE id = $1", locationId)
	return err
}
