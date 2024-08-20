package dio

import (
	"fmt"
	"io"
	"strings"

	"github.com/yznts/dsh/pkg/ddb"
	"go.kyoto.codes/zen/v3/jsonx"
	"go.kyoto.codes/zen/v3/slice"
)

// Csv is a writer that writes data as a multiple sql statements.
// It can write both schema and data, depending on the mode.
//
// This is kind of special writer,
// because it requires additional parameters to be set (mode and table).
// You have to keep attention on them.
type Sql struct {
	w io.Writer

	mode  string // one of "data", "schema"
	table string
}

// write wraps the io writer's Write method.
// If an error occurs, it panics.
// It's unexpected behavior in our case,
// so panic is necessary.
func (s *Sql) write(data []byte) {
	// Write and panic on error
	_, err := s.w.Write(data)
	if err != nil {
		panic(err)
	}
}

// Multi returns true if the writer supports multiple writes.
// Sql supports multiple writes (with multiple statements).
func (s *Sql) Multi() bool {
	return true
}

// WriteError usually outputs an error message.
// In our case we can't do that, so we panic.
func (s *Sql) WriteError(err error) {
	panic(fmt.Errorf("error while writing sql: %w", err))
}

// WriteData writes schema or data as a sql statement.
func (s *Sql) WriteData(data *ddb.Data) {

	// If we're writing schema, we need to write a CREATE TABLE statement
	// with taking data rows as column definitions.
	if s.mode == "schema" {
		// Convert data rows to column definitions
		col := strings.Join(slice.Map(data.Rows, func(row []any) string {
			return fmt.Sprintf("%s %s", row[0], row[1])
		}), ", \n")
		// Write the CREATE TABLE statement
		stm := fmt.Sprintf("CREATE TABLE %s (\n%s);\n\n", s.table, col)
		// Write the statement and return
		s.write([]byte(stm))
		return
	}

	// Otherwise, we're writing INSERT statement
	// with taking data rows as values.

	// First, let's write the INSERT statement.
	col := strings.Join(data.Cols, ", ")
	stm := fmt.Sprintf("INSERT INTO %s (%s) VALUES\n", s.table, col)
	s.write([]byte(stm))

	// And write data rows as values
	for i, row := range data.Rows {
		// If it's not the first row, write a comma and a newline
		if i != 0 {
			s.write([]byte(",\n"))
		}
		// Convert the row to a string slice.
		rowstr := strings.Join(slice.Map(row, func(val any) string {
			valstr := jsonx.String(val)
			if valstr[0] == '"' {
				valstr = fmt.Sprintf("'%s'", valstr[1:len(valstr)-1])
			}
			return valstr
		}), ", ")
		s.write([]byte(fmt.Sprintf("(%s)", rowstr)))
	}

	// Close the statement
	s.write([]byte(";\n\n"))
}

func (s *Sql) SetMode(mode string) {
	s.mode = mode
}

func (s *Sql) SetTable(table string) {
	s.table = table
}

func NewSql(w io.Writer) *Sql {
	return &Sql{w: w}
}
