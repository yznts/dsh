package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/yznts/dsh/pkg/dconf"
	"github.com/yznts/dsh/pkg/ddb"
	"github.com/yznts/dsh/pkg/dio"
	"go.kyoto.codes/zen/v3/slice"
)

// Tool flags
var (
	fdsn   = flag.String("dsn", "", "Database connection (can be set via DSN/DATABASE/DATABASE_URL env)")
	fall   = flag.Bool("all", false, "List all tables (including system)")
	fsql   = flag.Bool("sql", false, "Output in SQL format")
	fcsv   = flag.Bool("csv", false, "Output in CSV format")
	fjson  = flag.Bool("json", false, "Output in JSON format")
	fjsonl = flag.Bool("jsonl", false, "Output in JSON lines format")
)

// Tool usage / description
var (
	fusage = "[flags...] [table]"
	fdescr = "The dls utility lists tables/columns in the database."
)

// Database connection
var db ddb.Database

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
	stdout = dio.Open(os.Stdout, *fsql, *fcsv, *fjson, *fjsonl)
	stderr = dio.Open(os.Stderr, *fsql, *fcsv, *fjson, *fjsonl)

	// Resolve dsn and database connection
	dsn, err := dconf.GetDsn(*fdsn)
	dio.Error(stderr, err)
	db, err = ddb.Open(dsn)
	dio.Error(stderr, err)
	if db, iscloser := db.(io.Closer); iscloser {
		defer db.Close()
	}

	// Validate flags compatibility
	if *fall && *fsql {
		dio.Error(stderr, errors.New("flag -all is not compatible with -sql (export of system columns)"))
	}

	// If writer is SQL, we have a separate processing for it.
	if stdout, ok := stdout.(*dio.Sql); ok {
		// Determine tables we want to extract.
		// If no arguments, list all tables.
		// Otherwise, use provided table name.
		tables, err := db.QueryTables()
		dio.Error(stderr, err)
		if len(flag.Args()) > 0 {
			tables = slice.Filter(tables, func(t ddb.Table) bool {
				return t.Name == flag.Arg(0)
			})
		}

		// Filter system tables
		tables = slice.Filter(tables, func(t ddb.Table) bool {
			return !t.IsSystem
		})

		// Write schema for each table
		for _, table := range tables {
			// Get columns
			columns, err := db.QueryColumns(table.Name)
			dio.Error(stderr, err)
			// Set mode and table name
			stdout.SetMode("schema")
			stdout.SetTable(table.Name)
			// Write columns
			stdout.WriteData(&ddb.Data{
				Cols: []string{"COLUMN_NAME", "COLUMN_TYPE"},
				Rows: slice.Map(columns, func(c ddb.Column) []any {
					return []any{c.Name, c.Type}
				}),
			})
		}

		// Exit, we're done here
		return
	}

	// Otherwise, proceed with regular listing.

	// If no arguments, list tables.
	// Otherwise, list columns for provided table name.
	if len(flag.Args()) == 0 {
		// Get database tables
		tables, err := db.QueryTables()
		dio.Error(stderr, err)

		// Filter system tables
		if !*fall {
			tables = slice.Filter(tables, func(t ddb.Table) bool {
				return !t.IsSystem
			})
		}

		// If no schema, print 'N/A'
		if slice.All(tables, func(t ddb.Table) bool { return t.Schema == "" }) {
			tables = slice.Map(tables, func(t ddb.Table) ddb.Table {
				t.Schema = "N/A"
				return t
			})
		}

		// Write tables
		stdout.WriteData(&ddb.Data{
			Cols: []string{"TABLE_SCHEMA", "TABLE_NAME", "IS_SYSTEM"},
			Rows: slice.Map(tables, func(t ddb.Table) []any {
				return []any{t.Schema, t.Name, t.IsSystem}
			}),
		})
	} else {
		// Get database columns
		columns, err := db.QueryColumns(flag.Arg(0))
		dio.Error(stderr, err)

		// Write columns
		stdout.WriteData(&ddb.Data{
			Cols: []string{"COLUMN_NAME", "COLUMN_TYPE"},
			Rows: slice.Map(columns, func(c ddb.Column) []any {
				return []any{c.Name, c.Type}
			}),
		})
	}
}
