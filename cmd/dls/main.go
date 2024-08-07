package main

import (
	"database/sql"
	"flag"
	"net/url"
	"os"

	"github.com/yznts/dsh/pkg/ddb"
	"github.com/yznts/dsh/pkg/dio"
	"go.kyoto.codes/zen/v3/logic"
	"go.kyoto.codes/zen/v3/slice"
)

// Tool flags
var (
	fdsn   = flag.String("dsn", "", "Database connection (can be set via DSN/DATABASE/DATABASE_URL env)")
	fcsv   = flag.Bool("csv", false, "Output in CSV format")
	fjson  = flag.Bool("json", false, "Output in JSON format")
	fjsonl = flag.Bool("jsonl", false, "Output in JSON lines format")
)

// Tool usage / description
var (
	fusage = "[flags...] [table]"
	fdescr = "The dls utility lists tables (or table columns) in the database."
)

// Database connection
var (
	db    *sql.DB
	dbdsn *url.URL
)

// Output writers
var (
	stdout dio.Writer
	stderr dio.Writer
)

// Simplify assignments
var err error

func main() {
	// Provide usage
	flag.Usage = dio.Usage(fusage, fdescr)

	// Parse flags
	flag.Parse()

	// Resolve output writer
	stdout = dio.Open(os.Stdout, *fcsv, *fjson, *fjsonl)
	stderr = dio.Open(os.Stderr, *fcsv, *fjson, *fjsonl)

	// Resolve database connection
	db, dbdsn, err = ddb.Open(logic.Or(*fdsn,
		os.Getenv("DSN"),
		os.Getenv("DATABASE"),
		os.Getenv("DATABASE_URL")))
	dio.Error(stderr, err)

	// If no arguments, list tables.
	// Otherwise, list columns for provided table name.
	if len(flag.Args()) == 0 {
		// Get database tables
		tables, err := ddb.QueryTables(db, dbdsn.Scheme)
		dio.Error(stderr, err)

		// Write tables
		stdout.WriteTable(ddb.Data{
			Cols: []string{"TABLE_NAME"},
			Rows: slice.Map(tables, func(t ddb.Table) []any {
				return []any{t.Name}
			}),
		})
	} else {
		// Get database columns
		columns, err := ddb.QueryColumns(db, dbdsn.Scheme, flag.Arg(0))
		dio.Error(stderr, err)

		// Write columns
		stdout.WriteTable(ddb.Data{
			Cols: []string{"COLUMN_NAME", "COLUMN_TYPE"},
			Rows: slice.Map(columns, func(c ddb.Column) []any {
				return []any{c.Name, c.Type}
			}),
		})
	}
}
