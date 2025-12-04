//go:build dev || (!dev && !prod)

package config

import _ "embed"

//go:embed dev.yaml
var configData []byte

func init() {
	cfg, err := parseConfig(configData)
	if err != nil {
		panic(err)
	}

	cfg.PrintConfig()
}
