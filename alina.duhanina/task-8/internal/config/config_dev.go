//go:build dev || (!dev && !prod)

package config

import _ "embed"

func GetConfig() (*Config, error) {
	return parseConfig(devConfigData)
}
