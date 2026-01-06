package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (*AppConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("Load Config %s: no such file or directory", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Load Config: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("Loag Config: file %s: %w", path, err)
	}

	return &config, nil
}
