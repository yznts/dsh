package dio

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func Usage(usage, descr string) func() {
	return func() {
		// Let's limit the description width to 90 characters
		// to make it more readable.
		descr = lipgloss.NewStyle().Width(90).Render(descr)
		// Provide usage
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s %s\n\n", os.Args[0], usage)
		// Provide description
		fmt.Fprintf(flag.CommandLine.Output(), "%s \n\n", descr)
		// Provide flags,
		// this handled by flag package.
		flag.PrintDefaults()
	}
}
