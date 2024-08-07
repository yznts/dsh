package dio

import "io"

// Open returns a Writer based on the given flags.
func Open(
	w io.WriteCloser,
	csv, json, jsonl bool,
) Writer {
	switch {
	case csv:
		return NewCsv(w)
	case json:
		return NewJson(w)
	case jsonl:
		return NewJsonl(w)
	default:
		return NewPlain(w)
	}
}
