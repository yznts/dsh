package ddb

import (
	"database/sql"
	"database/sql/driver"
	"net/url"
	"reflect"

	"go.kyoto.codes/zen/v3/slice"
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
	// Execute the query.
	rows, err := c.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Get columns information.
	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	// Initialize the Data struct.
	// It holds both the columns and rows of the result.
	data := &Data{
		Cols: slice.Map(cols, func(c *sql.ColumnType) string { return c.Name() }),
	}
	// Define scan target row.
	// This is a slice of pointers,
	// so we need to copy values on each iteration.
	var scan []any
	for _, col := range cols {
		// Create a new pointer if corresponding column type.
		ptr := reflect.New(col.ScanType())
		// Append the pointer to the scan row
		scan = append(scan, ptr.Interface())
	}
	for rows.Next() {
		// Scan the row into prepared pointers
		err = rows.Scan(scan...)
		if err != nil {
			return nil, err
		}
		// Copy the values from the pointers to the Data struct
		var row []any
		for _, ptr := range scan {
			// If it's a nullable type, get the value
			if ptr, ok := ptr.(interface{ Value() (driver.Value, error) }); ok {
				val, _ := ptr.Value()
				row = append(row, val)
				continue
			}
			// Otherwise, get the value from the pointer
			row = append(row, reflect.ValueOf(ptr).Elem().Interface())
		}
		// Append the row to the Data holder
		data.Rows = append(data.Rows, row)
	}
	return data, nil
}
