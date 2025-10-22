package parseryaml

import (
	"fmt"
	"os"

	"github.com/TvoyBatyA12343/task-3/internal/config"
	"gopkg.in/yaml.v3"
)

func LoadConfig(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config config.Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
