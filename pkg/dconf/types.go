package dconf

// Configuration is the top-level configuration object
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
}
