package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/pkg/logger"
)

func CheckIDExists(table string, id int, db *pgxpool.Pool) (bool, error) {
	log := logger.NewLogger()
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", table)

	var exists bool
	err := db.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to check ID %d: %w", id, err)
		return false, err
	}

	return exists, nil
}
