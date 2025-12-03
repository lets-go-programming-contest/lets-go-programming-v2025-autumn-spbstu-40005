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

var currentEnv = "default"

func LoadConfig() (*Config, error) {
	var data []byte

	switch currentEnv {
	case "dev":
		data = devConfig
	case "prod":
		data = prodConfig
	default:
		data = defaultConfig
	}

	return parseConfig(data)
}

func parseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
