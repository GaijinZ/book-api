package pkg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: move it to postgres pkg
func CheckIDExists(table string, id int, db *pgxpool.Pool) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", table)

	var exists bool
	err := db.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking ID %d: %w", id, err)
	}

	return exists, nil
}
