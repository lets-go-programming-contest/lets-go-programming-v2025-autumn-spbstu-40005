package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
	"oleg.zholobov/task-3/internal/datamodels"
)

func LoadConfig(path string) (*datamodels.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config datamodels.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
