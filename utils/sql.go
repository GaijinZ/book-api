package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckIDExists(table string, id int, db *pgxpool.Pool) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", table)

	var exists bool
	err := db.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		errorMessage := fmt.Sprintf("Checking ID error %d: %s", id, err.Error())
		return false, fmt.Errorf(errorMessage)
	}

	return exists, nil
}
