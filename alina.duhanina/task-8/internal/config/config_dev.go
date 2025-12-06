//go:build dev
// +build dev

package config

func GetConfig() (*Config, error) {
	return parseConfig(devConfigData)
}
