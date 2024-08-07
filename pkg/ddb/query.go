/*
query.go contains wrappers for querying data from a database,
which must to be designed to work with any (supported) database.
*/
package ddb

import (
	"database/sql"
	"errors"
	"fmt"
)

// Database specific queries
var (
	sqlTablesSqlite   = "SELECT name FROM sqlite_master WHERE type='table';"
	sqlTablesPostgres = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';"

	sqlColumnsSqlite   = "SELECT name,type FROM PRAGMA_TABLE_INFO('%s')"
	sqlColumnsPostgres = "SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s'"
)

// QueryData queries the database and returns the result as a Data struct.
// The Data struct contains the columns and rows of the result.
// This format is designed to be used with our dio package.
func QueryData(db *sql.DB, query string) (Data, error) {
	rows, err := db.Query(query)
	if err != nil {
		return Data{}, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return Data{}, err
	}
	data := Data{
		Cols: cols,
	}
	for rows.Next() {
		var row []any
		for range cols {
			row = append(row, new(any))
		}
		err = rows.Scan(row...)
		if err != nil {
			return Data{}, err
		}
		var newRow []any
		for _, val := range row {
			newRow = append(newRow, *val.(*any))
		}
		data.Rows = append(data.Rows, newRow)
	}
	return data, nil
}

// QueryTables queries the database and returns the tables as a slice of Table structs.
func QueryTables(db *sql.DB, scheme string) ([]Table, error) {
	var (
		query  string
		tables []Table
	)
	switch ResolveScheme(scheme) {
	case "sqlite":
		query = sqlTablesSqlite
	case "postgres":
		query = sqlTablesPostgres
	default:
		return nil, errors.New("Unsupported database")
	}
	rows, err := db.Query(query)
	if err != nil {
		return tables, err
	}
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return tables, err
		}
		tables = append(tables, Table{Name: table})
	}
	return tables, nil
}

// QueryColumns gets the columns of a table and returns them as a slice of Column structs.
// The Column struct contains the data about the column (that can be provided by the database)
// which might be useful for further processing.
func QueryColumns(db *sql.DB, scheme string, table string) ([]Column, error) {
	var (
		query   string
		columns []Column
	)
	switch ResolveScheme(scheme) {
	case "sqlite", "sqlite3":
		query = fmt.Sprintf(sqlColumnsSqlite, table)
	case "postgres", "postgresql":
		query = fmt.Sprintf(sqlColumnsPostgres, table)
	default:
		return nil, errors.New("Unsupported database")
	}
	rows, err := db.Query(query)
	if err != nil {
		return []Column{}, err
	}
	for rows.Next() {
		var column, columnType string
		err = rows.Scan(&column, &columnType)
		if err != nil {
			return []Column{}, err
		}
		columns = append(columns, Column{
			Name: column,
			Type: columnType,
		})
	}
	return columns, nil
}
