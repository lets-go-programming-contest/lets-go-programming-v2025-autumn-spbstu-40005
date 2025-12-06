//go:build dev || (!dev && !prod)

package config

import _ "embed"

//go:embed dev.yaml
var configData []byte

func GetConfig() (*Config, error) {
	return parseConfig(configData)
}
