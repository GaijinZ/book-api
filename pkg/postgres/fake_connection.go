package postgres

import (
	"context"
	"library/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewFakeDB(ctx context.Context) (*DB, error) {
	log := utils.GetLogger(ctx)

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("err not expected while open a mock db, %v", err)
	}

	log.Infof("Postgres connected!")

	return &DB{
		DB:   db,
		mock: mock,
	}, nil
}

func (d *DB) GetMock() sqlmock.Sqlmock {
	return d.mock
}
