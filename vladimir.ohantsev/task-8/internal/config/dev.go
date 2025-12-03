//go:build dev || !prod

package config

import _ "embed"

//go:embed dev.yaml
var config []byte
