package ddb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.kyoto.codes/zen/v3/slice"
)

type Mysql struct {
	Connection
}

// QueryData is a method that queries the database
// with the given query and returns the result as a Data struct pointer.
//
// The Data struct contains the columns and rows of the result.
// Method is returning a pointer to avoid copying the Data struct,
// which might be large.
//
// MySQL driver doesn't make any type assertions on scan,
// so we need to utilize .ColumnTypes() information to get the correct types.
func (m *Mysql) QueryData(query string) (*Data, error) {
	// Execute the query.
	rows, err := m.Query(query)
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

func (m *Mysql) systemSchemas() []string {
	return []string{"mysql", "information_schema", "performance_schema", "sys"}
}

func (m *Mysql) QueryTables() ([]Table, error) {
	// Query the database for the tables
	data, err := m.QueryData("SELECT table_name,table_schema FROM information_schema.tables")
	if err != nil {
		return nil, err
	}
	// Convert the data to a slice of Table objects
	tables := slice.Map(data.Rows, func(r []any) Table {
		return Table{
			Name:   r[0].(string),
			Schema: r[1].(string),
		}
	})
	// Mark system tables
	tables = slice.Map(tables, func(t Table) Table {
		if slice.Contains(m.systemSchemas(), t.Schema) {
			t.IsSystem = true
		}
		return t
	})
	// Return
	return tables, nil
}

func (m *Mysql) QueryColumns(table string) ([]Column, error) {
	// Query the database for the columns
	data, err := m.QueryData(fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s'", table))
	if err != nil {
		return nil, err
	}
	// Convert the data to a slice of Column objects
	columns := slice.Map(data.Rows, func(r []any) Column {
		return Column{
			Name: r[0].(string),
			Type: r[1].(string),
		}
	})
	// Return
	return columns, nil
}

func (m *Mysql) QueryProcesses() ([]Process, error) {
	// Query the database for the currently running processes
	query := `
		SELECT id, time, user, db, info
		FROM information_schema.processlist
	`
	data, err := m.QueryData(query)
	if err != nil {
		return nil, err
	}

	// Convert the data to a slice of Process objects
	def := func(v any, def any) any {
		if v == nil {
			return def
		}
		return v
	}
	processes := slice.Map(data.Rows, func(r []any) Process {
		return Process{
			Pid:      int(def(r[0], 0).(uint64)),
			Duration: time.Duration(def(r[1], 0).(int32)) * time.Second,
			Username: def(r[2], "").(string),
			Database: def(r[3], "").(string),
			Query:    strings.Join(strings.Fields(def(r[4], "").(string)), " "),
		}
	})

	// Return the list of processes
	return processes, nil
}

func (m *Mysql) KillProcess(pid int, force bool) error {
	_, err := m.Exec(fmt.Sprintf("KILL %d", pid))
	return err
}
