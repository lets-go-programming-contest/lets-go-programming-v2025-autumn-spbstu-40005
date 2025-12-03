package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

//go:embed default.yaml
var defaultConfig []byte

//go:embed dev.yaml
var devConfig []byte

//go:embed prod.yaml
var prodConfig []byte

func loadDefaultConfig() (*Config, error) {
	return parseConfig(defaultConfig)
}

func loadDevConfig() (*Config, error) {
	return parseConfig(devConfig)
}

func loadProdConfig() (*Config, error) {
	return parseConfig(prodConfig)
}

func parseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
