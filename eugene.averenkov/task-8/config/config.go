package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string `yaml:"environment"`
	LogLevel string `yaml:"log_level"`
}

func (c Config) String() string {
	return fmt.Sprintf("%s %s", c.Env, c.LogLevel)
}

func parseYAML(data []byte) (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	return &cfg, nil
}
