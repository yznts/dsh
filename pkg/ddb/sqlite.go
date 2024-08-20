//go:build !daemon

package ddb

import (
	"errors"
	"fmt"

	"go.kyoto.codes/zen/v3/slice"
)

type Sqlite struct {
	Connection
}

func (s *Sqlite) systemTables() []string {
	return []string{"sqlite_master", "sqlite_sequence", "sqlite_stat1"}
}

func (s *Sqlite) QueryTables() ([]Table, error) {
	// Query the database for the tables
	data, err := s.QueryData("SELECT name,'' FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	// Convert the data to a slice of Table objects
	tables := slice.Map(data.Rows, func(r []any) Table {
		return Table{
			Name: r[0].(string),
		}
	})
	// SQLite doesn't include system tables into sqlite_master,
	// so we have to manually add them.
	tables = append(
		tables,
		slice.Map(s.systemTables(), func(t string) Table {
			return Table{Name: t, IsSystem: true}
		})...,
	)
	// Return
	return tables, nil
}

func (s *Sqlite) QueryColumns(table string) ([]Column, error) {
	// Query the database for the columns
	data, err := s.QueryData(fmt.Sprintf("SELECT name,type FROM PRAGMA_TABLE_INFO('%s')", table))
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

func (s *Sqlite) QueryProcesses() ([]Process, error) {
	return nil, errors.New("sqlite doesn't support process list query, use `lsof <file>` instead")
}

func (s *Sqlite) KillProcess(pid int, force bool) error {
	return errors.New("sqlite doesn't support process killing, use `lsof <file>` + `kill` instead")
}
