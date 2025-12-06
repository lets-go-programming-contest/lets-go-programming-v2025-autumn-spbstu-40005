//go:build prod
// +build prod

package config

func GetConfig() (*Config, error) {
	return parseConfig(prodConfigData)
}
