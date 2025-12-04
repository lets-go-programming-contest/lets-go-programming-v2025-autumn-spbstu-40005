//go:build dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devYaml []byte

func Load() (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(devYaml, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
