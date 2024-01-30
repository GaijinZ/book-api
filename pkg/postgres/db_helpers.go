package postgres

import (
	"database/sql"
	"fmt"
	"library/pkg/logger"
)

func CheckIDExists(table string, id int, db *sql.DB) (bool, error) {
	log := logger.NewLogger(2)

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", table)

	var exists bool
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to check ID %d: %v", id, err)
		return false, err
	}

	return exists, nil
}
