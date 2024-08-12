package dconf

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"go.kyoto.codes/zen/v3/slice"
	"gopkg.in/yaml.v3"
)

// Open reads a configuration file from the provided path.
// If no path is provided, it defaults to "$HOME/.dsh/config.*".
func Open(path ...string) (*Configuration, error) {
	// If no path is provided, we default to "$HOME/.dsh/config.*".
	// Configuration file might be in JSON or YAML format.
	if len(path) == 0 {
		home := os.Getenv("HOME")
		entries, err := os.ReadDir(home + "/.dsh")
		if err != nil {
			return nil, err
		}
		entries = slice.Filter(entries, func(e fs.DirEntry) bool {
			return slice.Contains([]string{"config.json", "config.yaml"}, e.Name())
		})
		if len(entries) > 0 {
			path = append(path, home+"/.dsh/"+entries[0].Name())
		}
	}

	// If no configuration file is found, it' s not an error.
	// We just return nil.
	if len(path) == 0 {
		return nil, nil
	}

	// Read the configuration file.
	confbts, err := os.ReadFile(path[0])
	if err != nil {
		return nil, err
	}

	// Unmarshal the configuration file,
	// depending on the file extension.
	conf := &Configuration{}
	switch filepath.Ext(path[0]) {
	case ".json":
		err = json.Unmarshal(confbts, conf)
	case ".yaml":
		err = yaml.Unmarshal(confbts, conf)
	}

	if err != nil {
		return nil, err
	}

	return conf, nil
}
