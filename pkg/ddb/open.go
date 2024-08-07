package ddb

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// Open opens a database connection based on the DSN.
// The DSN must be in the format scheme://user:password@host:port/database?param=value
func Open(dsn string) (*sql.DB, *url.URL, error) {
	// Parse the DSN
	dsnurl, err := url.Parse(dsn)
	if err != nil {
		return nil, nil, err
	}
	// Open the database, depending on the scheme
	var db *sql.DB
	switch ResolveScheme(dsnurl.Scheme) {
	case "sqlite":
		// To open a SQLite database, we need to remove the scheme and leading slashes
		_dsnurl, _ := url.Parse(dsn)
		_dsnurl.Scheme = ""
		_dsnurlstr := strings.ReplaceAll(_dsnurl.String(), "//", "")
		db, err = sql.Open("sqlite", _dsnurlstr)
	case "postgres":
		db, err = sql.Open("pgx", dsn)
	default:
		err = errors.New("Empty DSN or unsupported database")
	}
	// Return the database connection
	return db, dsnurl, err
}
