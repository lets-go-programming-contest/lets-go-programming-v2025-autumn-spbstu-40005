package main

import (
	"flag"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/config"
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

	config, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	valutes, err := converter.ParseXMLFile(config.InputFile)
	if err != nil {
		panic(err)
	}

	converter.SortValutes(valutes)

	if err := converter.WriteToJSON(valutes, config.OutputFile); err != nil {
		panic(err)
	}
}
