package dconf

// Default configuration,
// loaded from $HOME/.dsh/config.*.
// If something goes wrong, DefaultErr will be set.
// Please note, even if DefaultErr is nil,
// Default also can be nil (no configuration found).
var (
	Default    *Configuration
	DefaultErr error
)

func init() {
	Default, DefaultErr = Open()
}
