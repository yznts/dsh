package ddb

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

func Open(dsn string) (Database, error) {
	// Validate and parse dsn
	if dsn == "" {
		return nil, errors.New("empty DSN")
	}
	dsnurl, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	// Resolve connection and actual scheme, depending on the DSN scheme
	switch dsnurl.Scheme {

	case "sqlite", "sqlite3":
		// To open a SQLite database, we need to remove the scheme and leading slashes
		_dsnurl, _ := url.Parse(dsn)
		_dsnurl.Scheme = ""
		_dsnurlstr := strings.ReplaceAll(_dsnurl.String(), "//", "")
		// Open sql database connection
		sqldb, err := sql.Open("sqlite", _dsnurlstr)
		if err != nil {
			return nil, err
		}
		// Compose the database object
		return &Sqlite{
			Connection: Connection{
				DB:     sqldb,
				DSN:    dsnurl,
				Scheme: "sqlite",
			},
		}, nil

	case "postgres":
		// Open sql database connection
		sqldb, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		// Compose the database object
		return &Postgres{
			Connection: Connection{
				DB:     sqldb,
				DSN:    dsnurl,
				Scheme: "postgres",
			},
		}, nil

	default:
		return nil, errors.New("unsupported database")
	}
}
