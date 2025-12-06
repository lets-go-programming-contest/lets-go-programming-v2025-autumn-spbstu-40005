//go:build prod

package config

import _ "embed"

func GetConfig() (*Config, error) {
	return parseConfig(prodConfigData)
}
