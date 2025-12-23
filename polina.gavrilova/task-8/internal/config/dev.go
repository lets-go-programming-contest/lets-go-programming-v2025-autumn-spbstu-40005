//go:build !prod
// +build !prod

package config

import _ "embed"

//go:embed dev.yaml
var configFile []byte
