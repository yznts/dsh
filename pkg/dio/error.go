package dio

import (
	"errors"
	"os"
)

// Error checks if the error presents,
// writes the error to the writer and exits the program with non-zero code.
// Override allows to provide a custom error message.
func Error(w Writer, err error, override ...string) {
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			panic(err)
		}
		if len(override) > 0 {
			w.WriteError(errors.New(override[0]))
		} else {
			w.WriteError(err)
		}
		os.Exit(1)
	}
}
