package yamlparsing

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"gordey.shapkov/task-3/internal/config"
)

func ParseYAMLFile(path string) (*config.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	cfg := &config.Config{InputFile: "", OutputFile: ""}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	return cfg, nil
}
