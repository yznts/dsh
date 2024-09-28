package dconf

import "os"

// Default configuration,
// loaded from $HOME/.dsh/config.*.
// If something goes wrong, DefaultErr will be set.
// Please note, even if DefaultErr is nil,
// Default also can be nil (no configuration found).
var (
	Default    *Configuration
	DefaultErr error
)

// init loads the default configuration.
func init() {
	Default, DefaultErr = OpenDefault()
}

// OpenDefault reads a configuration file from a predefined paths.
// It tries to resolve from multiple common locations:
// - "$HOME/.dsh/config.{json,yaml}"
// - "$HOME/.config/dsh/config.{json,yaml}"
func OpenDefault() (*Configuration, error) {
	// Defmine common locations
	home := os.Getenv("HOME")
	locations := []string{
		home + "/.dsh/config.json",
		home + "/.dsh/config.yaml",
		home + "/.config/dsh/config.json",
		home + "/.config/dsh/config.yaml",
	}
	// We are taking the first configuration file found
	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return Open(location)
		}
	}
	// If no configuration file is found, it' s not an error.
	// We just return nil.
	return nil, nil
}
