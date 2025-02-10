package dat

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Location struct {
	Id         int
	Created    time.Time
	Name       string
	Address    string
	Longtitude float64
	Latitude   float64
}

type Locations interface {
	InsertLocationIntoDatabase(loc Location) error
	GetLocationDetailsByID(locationId int) (*Location, error)
	GetAllLocationDetails() (*[]Location, error)
	UpdateLocationDetails(int, Location) error
	DeleteLocationFromDatabase(int) error
}

func (l *Location) ScanRow(r Row) error {
	return r.Scan(
		&l.Id,
		&l.Name,
		&l.Address,
		&l.Longtitude,
		&l.Latitude,
	)
}

func (c *Context) InsertLocationIntoDatabase(loc Location) error {
	query := "INSERT INTO locations (name, address, longtitude, latitude) VALUES ($1, $2, $3, $4)"

	tx, err := c.ConnectionPool.BeginTx(context.TODO(), pgx.TxOptions{})
	if err != nil {
		c.Logger.Error("Failed to begin transaction", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
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
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
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
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
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

func UpdateLocationDetails(int, Location) error {
	return nil
}

func DeleteLocationFromDatabase(int) error {
	return nil
}
