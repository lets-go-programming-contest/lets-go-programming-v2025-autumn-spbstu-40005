//go:build prod

package config

//go:embed prod.yaml
var prodConfig []byte

func LoadConfig() (*Config, error) {
	return loadConfig(prodConfig)
}
