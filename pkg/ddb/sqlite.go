package ddb

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

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

func OpenSqlite(dsn string) (*Sqlite, error) {
	// Parse the DSN
	dsnurl, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	// To open a SQLite database, we need to remove the scheme and leading slashes
	_dsnurl, _ := url.Parse(dsn)
	_dsnurl.Scheme = ""
	_dsnurlstr := strings.ReplaceAll(_dsnurl.String(), "//", "")
	sqldb, err := sql.Open("sqlite", _dsnurlstr)
	if err != nil {
		return nil, err
	}
	// Return
	return &Sqlite{
		Connection: Connection{
			DB:     sqldb,
			DSN:    dsnurl,
			Scheme: "sqlite",
		},
	}, nil
}
