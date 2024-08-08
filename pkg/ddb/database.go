package ddb

// Database interface summarizes the methods
// that our utilities are going to use to interact with databases.
//
// Most of the methods are database-specific and must be implemented
// by the database-specific struct.
// On the other hand, database-agnostic methods might be implemented
// on Connection struct, which nested into each database-specific struct.
type Database interface {
	QueryData(query string) (*Data, error) // Return a pointer because data amount might be large
	QueryTables() ([]Table, error)
	QueryColumns(table string) ([]Column, error)
}
