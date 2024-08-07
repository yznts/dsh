package dio

import "os"

// Error checks if the error presents,
// writes the error to the writer and exits the program with non-zero code.
func Error(w Writer, err error) {
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			panic(err)
		}
		w.WriteError(err)
		os.Exit(1)
	}
}
