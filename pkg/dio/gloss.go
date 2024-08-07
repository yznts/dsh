package dio

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/yznts/dsh/pkg/ddb"
	"go.kyoto.codes/zen/v3/slice"
)

// Gloss is a writer that writes a formatted output,
// like a table or styled error/warn messages.
// Uses lipgloss for styling,
// that's why it's called Gloss.
type Gloss struct {
	w io.WriteCloser
}

// write wraps the io writer's Write method.
// If an error occurs, it panics.
// It's unexpected behavior in our case,
// so panic is necessary.
func (g *Gloss) write(data []byte) {
	_, err := g.w.Write(data)
	if err != nil {
		panic(err)
	}
}

// Multi returns true if the writer supports multiple writes.
// Gloss does not support multiple writes,
// because it outputs in a formatted way that cannot be appended to (i.e. closed table).
func (g *Gloss) Multi() bool {
	return false
}

func (g *Gloss) WriteError(err error) {
	msg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f66f81")).
		Bold(true).
		Render(fmt.Sprintf("error occured: %s", err.Error()))
	g.write([]byte(msg + "\n"))
	// No need to close writer, because it's just an error message.
	// We can write more data after that.
}

func (g *Gloss) WriteTable(data ddb.Data) {
	// Transform rows to string
	rowsstr := slice.Map(data.Rows, func(v []any) []string {
		return slice.Map(v, func(v any) string {
			return fmt.Sprintf("%v", v)
		})
	})
	// Create table
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true).Padding(0, 2)
			} else {
				return lipgloss.NewStyle().MaxHeight(5).MaxWidth(80).Padding(0, 2)
			}
		}).
		Headers(data.Cols...).
		Rows(rowsstr...)
	// Write table
	g.write([]byte(t.String() + "\n"))
	// Close writer.
	// After the table is written, it cannot be appended to.
	// If someone will try to write once more, it will panic.
	if err := g.w.Close(); err != nil {
		panic(err)
	}
}

func (g *Gloss) WriteWarning(msg string) {
	_msg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f6ef6f")).
		Bold(true).
		Render(fmt.Sprintf("warning: %s", msg))
	g.w.Write([]byte(_msg + "\n"))
	// No need to close writer, because it's just a warning message.
	// We can write more data after that.
}

func NewPlain(w io.WriteCloser) *Gloss {
	return &Gloss{w: w}
}
