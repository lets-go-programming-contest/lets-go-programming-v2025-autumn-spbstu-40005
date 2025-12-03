//go:build dev

package config

func LoadConfig() (*Config, error) {
	return loadDevConfig()
}
