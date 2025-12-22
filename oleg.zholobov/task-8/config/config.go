package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func (c *Config) String() string {
	return fmt.Sprintf("%s %s", c.Environment, c.LogLevel)
}

func GetConfig(configData []byte) (*Config, error) {
	cfg := &Config{
		Environment: "",
		LogLevel: "",
	}

	err := yaml.Unmarshal(configData, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
