package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func ParseConfig() (*Config, error) {
	var config Config

	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parsing config: %w", err)
	}

	return &config, nil
}

func (c *Config) PrintConfig() {
	fmt.Print(c.Environment, " ", c.LogLevel)
}
