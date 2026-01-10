package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env string `yaml:"environment"`
	Log string `yaml:"log_level"`
}

func GetConfig() (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(cfgData, &c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &c, nil
}
