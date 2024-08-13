package ddb

import (
	"fmt"
	"strings"
	"time"

	"go.kyoto.codes/zen/v3/slice"
)

type Mysql struct {
	Connection
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
