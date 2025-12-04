//go:build !dev

package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodYaml []byte

func Load() (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(prodYaml, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
