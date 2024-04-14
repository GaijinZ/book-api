package postgres

import (
	"context"
	"database/sql"
	"library/pkg/utils"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	DriverName      string
	DataSourceName  string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type DB struct {
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func NewDB(ctx context.Context, configDB DBConfig) (*DB, error) {
	log := utils.GetLogger(ctx)

	db, err := sql.Open(configDB.DriverName, configDB.DataSourceName)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", configDB.DriverName)
		return nil, err
	}

	db.SetMaxOpenConns(configDB.MaxOpenConns)
	db.SetMaxIdleConns(configDB.MaxIdleConns)
	db.SetConnMaxLifetime(configDB.ConnMaxLifetime)

	return &DB{DB: db}, nil
}

func (d *DB) GetDB() *sql.DB {
	return d.DB
}

func (d *DB) Ping() error {
	return d.DB.Ping()
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func (d *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	return d.DB.Conn(ctx)
}
