package dconf

import (
	"encoding/json"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Open reads a configuration file from the given path.
func Open(path string) (*Configuration, error) {
	// Read the configuration file.
	confbts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the configuration file,
	// depending on the file extension.
	conf := &Configuration{}
	switch filepath.Ext(path) {
	case ".json":
		err = json.Unmarshal(confbts, conf)
	case ".yaml":
		err = yaml.Unmarshal(confbts, conf)
	}

	// Return
	return conf, err
}
