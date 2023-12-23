package postgres

import (
	"fmt"
	"library/pkg/logger"
)

func CheckIDExists(table string, id int, db DB) (bool, error) {
	log := logger.NewLogger(2)

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", table)

	var exists bool
	err := db.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to check ID %d: %w", id, err)
		return false, err
	}

	return exists, nil
}
