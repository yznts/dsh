package dconf

import (
	"errors"
	"net/url"
	"os"

	"go.kyoto.codes/zen/v3/logic"
)

// GetDsn is a common dsn resolver.
// It tries to resolve provided dsn string to the actual dsn.
// Sometimes it might be empty (expecting env resolving),
// sometimes it's a name from configuration file (expecting conf resolving).
//
// It tries to resolve dsn in the following order:
// - if dsn is empty, it tries to resolve it from environment variable
// - if dsn schema not found, it tries to resolve actual dsn from configuration file
// - if dsn schema found, it returns dsn as is
func GetDsn(dsn string) (string, error) {
	// If dsn is empty, try to resolve it from environment variable
	dsn = logic.Or(dsn,
		os.Getenv("DSN"),
		os.Getenv("DATABASE"),
		os.Getenv("DATABASE_URL"))
	// If it's still empty, return an error
	if dsn == "" {
		return "", errors.New("dsn is empty")
	}
	// Parse dsn
	dsnurl, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}
	// If dsn schema not found, try to resolve actual dsn from default configuration file
	if dsnurl.Scheme == "" {
		// If connection configuration found, go through the same process.
		// Otherwise, it's dsn error.
		if con, ok := Default.GetConnection(dsnurl.Path); ok {
			return GetDsn(con.Conn)
		} else {
			return "", errors.New("dsn schema not found")
		}
	}
	// Return dsn as is
	return dsn, nil
}
