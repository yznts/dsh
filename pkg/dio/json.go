package dio

import (
	"io"

	"github.com/yznts/dsh/pkg/ddb"
	"go.kyoto.codes/zen/v3/jsonx"
)

// Json is a writer that writes a single json object.
type Json struct {
	w io.WriteCloser
}

// write wraps the io writer's Write method.
// If an error occurs, it panics.
// It's unexpected behavior in our case,
// so panic is necessary.
func (j *Json) write(data []byte) {
	if _, err := j.w.Write(data); err != nil {
		panic(err)
	}
	if err := j.w.Close(); err != nil {
		panic(err)
	}
}

// Multi returns true if the writer supports multiple writes.
// Json does not support multiple writes,
// because it must output a single JSON object (unlike JSONL).
func (j *Json) Multi() bool {
	return false
}

func (j *Json) WriteError(err error) {
	errmap := map[string]any{"error": err.Error()}
	j.write(jsonx.Bytes(errmap))
}

func (j *Json) WriteTable(data ddb.Data) {
	j.write(jsonx.Bytes(map[string]any{
		"cols": data.Cols,
		"rows": data.Rows,
	}))
}

func NewJson(w io.WriteCloser) *Json {
	return &Json{w: w}
}
