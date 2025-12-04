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
	var config Config

	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parsing config: %w", err)
	}

	return &config, nil
}

func (c *Config) PrintConfig() {
	fmt.Print(c.Environment, " ", c.LogLevel)
}
