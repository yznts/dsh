package main

import (
	"database/sql"
	"flag"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/yznts/dsh/pkg/ddb"
	"github.com/yznts/dsh/pkg/dio"
	"go.kyoto.codes/zen/v3/logic"
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
	fusage = "[flags...] sql"
	fdescr = "The dsql utility executes SQL query and writes the result to the standard output in desired format. " +
		"It designed to be simple, therefore edge cases handling isn't included, like trying to query large tables in a formatted way. \n\n" +
		"The query can be provided as argument or piped from another command (STDIN). "
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

	// Extract sql query from arguments
	query := strings.Join(flag.Args(), " ")
	// If no query provided, read from STDIN
	if query == "" {
		querybts, err := io.ReadAll(os.Stdin)
		dio.Error(stderr, err)
		query = string(querybts)
	}

	// Execute the query
	data, err := ddb.QueryData(db, query)
	dio.Error(stderr, err)

	// Write the result
	stdout.WriteTable(data)
}
