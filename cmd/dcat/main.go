package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/yznts/dsh/pkg/ddb"
	"github.com/yznts/dsh/pkg/dio"
	"go.kyoto.codes/zen/v3/logic"
	"go.kyoto.codes/zen/v3/slice"
)

// Tool flags
var (
	fdsn   = flag.String("dsn", "", "Database connection (can be set via DSN/DATABASE/DATABASE_URL env)")
	fcsv   = flag.Bool("csv", false, "Output in CSV format")
	fjsonl = flag.Bool("jsonl", false, "Output in JSON lines format")
	fwhere = flag.String("where", "", "WHERE clause")
)

// Tool usage / description
var (
	fusage = "[flags...] table"
	fdescr = "The dcat utility reads table data and writes it to the standard output in desired format. " +
		"Because of chunked data fetching, output options might be limited. " +
		"Utility tries to avoid accumulating data in the memory. " +
		"If dcat output options are not enough and memory usage is not a concern, consider using dsql instead."
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

	// Resolve output writer.
	//
	// We are using only multiline writers here
	// because we are going to query the database in chunks
	// and write the result in chunks as well.
	// Otherwise, we will have to store the whole result in memory.
	//
	// The only exception is gloss writer.
	// We will limit the output to 1k rows in that case.
	stdout = dio.Open(os.Stdout, *fcsv, false, *fjsonl)
	stderr = dio.Open(os.Stderr, *fcsv, false, *fjsonl)

	// Determine if the output format supports multiple writes.
	// Otherwise, we are limiting the output to 1k rows with a warning.
	limited := false
	if !stdout.Multi() {
		limited = true
		// Let's also check if the output format is gloss.
		// Other formats are not supposed to be used without multiple writes.
		if _, gloss := stdout.(*dio.Gloss); !gloss {
			dio.Error(stderr, errors.New("output format does not support multiple writes"))
		}
	}

	// Resolve database connection
	db, dbdsn, err = ddb.Open(logic.Or(*fdsn,
		os.Getenv("DSN"),
		os.Getenv("DATABASE"),
		os.Getenv("DATABASE_URL")))
	dio.Error(stderr, err)

	// Extract table name from arguments
	table := flag.Arg(0)
	if table == "" {
		dio.Error(stderr, errors.New("missing table name"))
	}

	// Get rows count
	count, err := ddb.Count(db, table, *fwhere)
	dio.Error(stderr, err)

	// If the total rows count is less than 1k, we can get back to non-limited mode.
	if count < 1000 {
		limited = false
	}

	// Make chunks
	chunks := slice.Chunks(slice.Range(0, count), 1000)

	// We need only first item from each chunk as an offset.
	offsets := slice.Map(chunks, func(chunk []int) int { return chunk[0] })
	// If we are limited, we need only first chunk.
	// Also, we need to warn the user about it.
	if limited {
		offsets = offsets[:1]
		if stdout, warner := stdout.(dio.WarningWriter); warner {
			stdout.WriteWarning("output is limited to 1k rows")
		}
	}

	// Iterate over chunks and query the database
	for _, offset := range offsets {
		// Compose limit/offset query with WHERE clause
		query := &strings.Builder{}
		query.WriteString(fmt.Sprintf("SELECT * FROM %s ", table))
		if *fwhere != "" {
			query.WriteString(fmt.Sprintf("WHERE %s ", *fwhere))
		}
		query.WriteString(fmt.Sprintf("LIMIT 1000 OFFSET %d", offset))
		// Execute query
		data, err := ddb.QueryData(db, query.String())
		dio.Error(stderr, err)
		// Don't collect the data and just write it to the output,
		// because we don't want to keep it in memory.
		// That's why we are requiring closable writers here.
		stdout.WriteTable(data)
	}
}
