//go:build dev

package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devYaml []byte

func init() {
	yaml.Unmarshal(devYaml, &Cfg)
}
