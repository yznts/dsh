package ddb

import (
	"database/sql"
)

func Count(db *sql.DB, table string, where ...string) (int, error) {
	query := "SELECT COUNT(*) FROM " + table
	if len(where) > 0 && where[0] != "" {
		query += " WHERE " + where[0]
	}
	var count int
	err := db.QueryRow(query).Scan(&count)
	return count, err
}
