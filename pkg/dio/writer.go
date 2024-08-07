package dio

import (
	"github.com/yznts/dsh/pkg/ddb"
)

// Writer is an interface that must be implemented by all writers.
// It provides a common interface for all tools to write data/errors/etc.
type Writer interface {
	Multi() bool // Multi returns true if the writer supports multiple writes.
	WriteError(error)
	WriteTable(ddb.Data)
}

// WarningWriter is an optional interface that can be implemented by writers.
// It allows writers to report warnings.
type WarningWriter interface {
	WriteWarning(string)
}
