//go:build prod

package config

func LoadConfig() (*Config, error) {
	return loadProdConfig()
}
