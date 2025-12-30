//go:build prod || (!dev && !prod)

package config

import _ "embed"

//go:embed prod.yaml
var configData []byte

func GetConfig() (*Config, error) {
	return parseConfig(configData)
}
