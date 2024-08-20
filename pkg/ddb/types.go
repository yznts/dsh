package ddb

import "time"

// Database interface summarizes the methods
// that our utilities are going to use to interact with databases.
//
// Most of the methods are database-specific and must be implemented
// by the database-specific struct.
// On the other hand, database-agnostic methods might be implemented
// on Connection struct, which nested into each database-specific struct.
type Database interface {
	// Data queries
	QueryData(query string) (*Data, error) // Return a pointer because data amount might be large

	// Schema queries
	QueryTables() ([]Table, error)
	QueryColumns(table string) ([]Column, error)

	// Process queries
	QueryProcesses() ([]Process, error)
	KillProcess(pid int, force bool) error
}

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
