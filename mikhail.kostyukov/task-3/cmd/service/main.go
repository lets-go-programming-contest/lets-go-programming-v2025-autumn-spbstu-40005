package main

import (
	"flag"
	"fmt"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/converter"
)

var (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %v\n", err)
		}
	}()

	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	if err := converter.Run(*configPath); err != nil {
		fmt.Printf("converter.Run: %v\n", err)
	}
}
