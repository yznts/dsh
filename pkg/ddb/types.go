package ddb

import "time"

// Data holds query results.
// Columns and rows are stored separately instead of using maps,
// so we can minimize memory usage and output.
type Data struct {
	Cols []string
	Rows [][]any
}

// Table holds table meta information,
// not the actual data.
type Table struct {
	Schema   string
	Name     string
	IsSystem bool // Indicates whether it's a system table
}

// Column holds column meta information.
type Column struct {
	Name string
	Type string
}

type Process struct {
	Pid      int
	Duration time.Duration
	Username string
	Database string
	Query    string
}
