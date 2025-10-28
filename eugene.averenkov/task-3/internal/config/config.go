package config

import (
    "fmt"
    "os"
    "gopkg.in/yaml.v3"
)

type Config struct {
    InputFile  string `yaml:"input-file"`
    OutputFile string `yaml:"output-file"`
}

func Load(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
    }

    if cfg.InputFile == "" {
        return nil, fmt.Errorf("input-file is required")
    }
    if cfg.OutputFile == "" {
        return nil, fmt.Errorf("output-file is required")
    }

    return &cfg, nil
}
