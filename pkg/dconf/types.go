package dconf

import (
	"net/url"
	"strconv"

	"go.kyoto.codes/zen/v3/slice"
)

// Configuration is the top-level configuration object.
// Please note, it holds only the data stored in the configuration file
// and doesn't take into account the environment variables, args, etc.
type Configuration struct {
	Connections []Connection `json:"connections" yaml:"connections"`
}

func (c *Configuration) GetConnection(name string) (Connection, bool) {
	for _, conn := range c.Connections {
		if conn.Name == name {
			return conn, true
		}
	}
	return Connection{}, false
}

// Connection is a connection object
type Connection struct {
	Name string `json:"name" yaml:"name"`

	// Conn is the raw connection string, DSN.
	// It will be passed to the driver as-is.
	Conn string `json:"conn" yaml:"conn"`

	// As an alternative,
	// we're giving the ability to specify connection parameters separately.
	// This approach might be more convenient for reading and modifying.
	Type    string `json:"type" yaml:"type"`
	Host    string `json:"host" yaml:"host"`
	Port    int    `json:"port" yaml:"port"`
	User    string `json:"user" yaml:"user"`
	Pass    string `json:"pass" yaml:"pass"`
	DB      string `json:"db" yaml:"db"`
	SslMode string `json:"ssl_mode" yaml:"ssl_mode"`
	SslCert string `json:"ssl_cert" yaml:"ssl_cert"`
	SslKey  string `json:"ssl_key" yaml:"ssl_key"`
	SslCa   string `json:"ssl_ca" yaml:"ssl_ca"`
}

// GetConn returns the connection string.
// If the Conn field is set, it will be returned as is.
// Otherwise, the connection string will be built from the separate fields.
func (c *Connection) GetConn() string {
	// If the Conn field is set, return it as is.
	if c.Conn != "" {
		return c.Conn
	}
	// Otherwise, build the DSN from the separate fields.
	dsn := &url.URL{}
	dsn.Scheme = c.Type
	dsn.User = url.UserPassword(c.User, c.Pass)
	dsn.Host = c.Host
	if c.Port != 0 {
		dsn.Host += ":" + strconv.Itoa(c.Port)
	}
	dsn.Path = c.DB
	q := dsn.Query()
	// Query parameters format may vary depending on the driver.
	switch {
	case slice.Contains([]string{"postgres", "postgresql"}, c.Type):
		q.Add("sslmode", c.SslMode)
		q.Add("sslcert", c.SslCert)
		q.Add("sslkey", c.SslKey)
		q.Add("sslrootcert", c.SslCa)
	case slice.Contains([]string{"mysql"}, c.Type):
		q.Add("ssl_mode", c.SslMode)
		q.Add("ssl_cert", c.SslCert)
		q.Add("ssl_key", c.SslKey)
		q.Add("ssl_ca", c.SslCa)
	}
	dsn.RawQuery = q.Encode()
	return dsn.String()
}
