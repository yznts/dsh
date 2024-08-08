package ddb

import (
	"database/sql"
	"net/url"
)

// Connection is a wrapper around sql.DB that also stores the DSN and scheme.
// Also, it holds database-agnostic methods.
type Connection struct {
	*sql.DB

	DSN    *url.URL
	Scheme string
}

// QueryData is a database-agnostic method that queries the database
// with the given query and returns the result as a Data struct pointer.
//
// The Data struct contains the columns and rows of the result.
//
// Method is returning a pointer to avoid copying the Data struct,
// which might be large.
func (c *Connection) QueryData(query string) (*Data, error) {
	rows, err := c.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	data := &Data{
		Cols: cols,
	}
	for rows.Next() {
		var row []any
		for range cols {
			row = append(row, new(any))
		}
		err = rows.Scan(row...)
		if err != nil {
			return nil, err
		}
		var newRow []any
		for _, val := range row {
			newRow = append(newRow, *val.(*any))
		}
		data.Rows = append(data.Rows, newRow)
	}
	return data, nil
}
