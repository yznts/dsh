package ddb

// Data holds query results.
// Columns and rows are stored separately instead of using maps,
// so we can minimize memory usage and output.
type Data struct {
	Cols []string
	Rows [][]any
}

type Table struct {
	Name string
}

type Column struct {
	Name string
	Type string
}
