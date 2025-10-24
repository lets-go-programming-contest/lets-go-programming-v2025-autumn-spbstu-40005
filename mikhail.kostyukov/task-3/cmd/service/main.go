package main

import (
	"flag"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/converter"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	if err := converter.Run(*configPath); err != nil {
		panic(err)
	}
}
