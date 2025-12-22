//go:build dev
// +build dev

package config

import _ "embed"

//go:embed dev.yaml
var configData []byte

func LoadConfig() (*Config, error) {
	return GetConfig(configData)
}
