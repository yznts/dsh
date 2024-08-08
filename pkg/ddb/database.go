package ddb

type Database interface {
	QueryData(query string) (*Data, error) // Return a pointer because data amount might be large
	QueryTables() ([]Table, error)
	QueryColumns(table string) ([]Column, error)
}
