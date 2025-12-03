//go:build !dev && !prod

package config

func init() {
	currentEnv = "default"
}
