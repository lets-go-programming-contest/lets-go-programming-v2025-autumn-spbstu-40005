//go:build !dev

package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodYaml []byte

func init() {
	yaml.Unmarshal(prodYaml, &Cfg)
}
