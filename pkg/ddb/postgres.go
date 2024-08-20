//go:build !daemon

package ddb

import (
	"fmt"
	"strings"
	"time"

	"go.kyoto.codes/zen/v3/slice"
)

type Postgres struct {
	Connection
}

func (p *Postgres) systemSchemas() []string {
	return []string{"pg_catalog", "information_schema"}
}

func (p *Postgres) QueryTables() ([]Table, error) {
	// Query the database for the tables
	data, err := p.QueryData("SELECT table_name,table_schema FROM information_schema.tables")
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
		if slice.Contains(p.systemSchemas(), t.Schema) {
			t.IsSystem = true
		}
		return t
	})
	// Return
	return tables, nil
}

func (p *Postgres) QueryColumns(table string) ([]Column, error) {
	// Query the database for the columns
	data, err := p.QueryData(fmt.Sprintf("SELECT column_name,data_type FROM information_schema.columns WHERE table_name = '%s'", table))
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

func (p *Postgres) QueryProcesses() ([]Process, error) {
	// Query the database for the currently running processes
	query := `
		SELECT
			pid,
			date_part('epoch', now() - pg_stat_activity.query_start) AS duration,
			usename,
			datname,
			query
		FROM
			pg_stat_activity
	`
	data, err := p.QueryData(query)
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
			Pid:      int(def(r[0], 0).(int64)),
			Duration: time.Duration(def(r[1], 0.0).(float64)) * time.Second,
			Username: def(r[2], "").(string),
			Database: def(r[3], "").(string),
			Query:    strings.Join(strings.Fields(def(r[4], "").(string)), " "),
		}
	})

	// Return the list of processes
	return processes, nil
}

func (p *Postgres) KillProcess(pid int, force bool) error {
	if !force {
		_, err := p.Exec(fmt.Sprintf("SELECT pg_cancel_backend(%d)", pid))
		return err
	} else {
		_, err := p.Exec(fmt.Sprintf("SELECT pg_terminate_backend(%d)", pid))
		return err
	}
}
