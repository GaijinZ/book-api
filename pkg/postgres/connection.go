package postgres

import (
	"context"
	"database/sql"
	"library/pkg/utils"
	"time"

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
	DB *sql.DB
}

func NewDB(ctx context.Context, configDB DBConfig) (*DB, error) {
	log := utils.GetLogger(ctx)

	db, err := sql.Open(configDB.DriverName, configDB.DataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(configDB.MaxOpenConns)
	db.SetMaxIdleConns(configDB.MaxIdleConns)
	db.SetConnMaxLifetime(configDB.ConnMaxLifetime)

	log.Infof("Postgres connected!")

	return &DB{DB: db}, nil
}

func (d *DB) Ping() error {
	return d.DB.Ping()
}

func (d *DB) Close() error {
	return d.DB.Close()
}

func (d *DB) GetConnection(ctx context.Context) (*sql.Conn, error) {
	return d.DB.Conn(ctx)
}
