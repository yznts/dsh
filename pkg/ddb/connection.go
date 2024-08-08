package ddb

import (
	"database/sql"
	"net/url"
)

type Connection struct {
	*sql.DB

	DSN    *url.URL
	Scheme string
}

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
