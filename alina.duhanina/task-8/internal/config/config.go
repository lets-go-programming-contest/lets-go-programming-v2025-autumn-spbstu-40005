package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func parseConfig(data []byte) (*Config, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty config data")
	}

	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parsing config: %w", err)
	}

	if config.Environment == "" {
		return nil, fmt.Errorf("environment field is empty")
	}
	if config.LogLevel == "" {
		return nil, fmt.Errorf("log_level field is empty")
	}
	return &config, nil
}

func (c *Config) PrintConfig() {
	fmt.Print(c.Environment, " ", c.LogLevel)
}
