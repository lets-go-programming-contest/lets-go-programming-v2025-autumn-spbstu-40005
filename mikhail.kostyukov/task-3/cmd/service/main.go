package main

import (
	"flag"
	"sort"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/config"
	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/converter"
	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/model"
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

	var rates model.ValCurs
	if err := converter.ParseXMLFile(config.InputFile, &rates); err != nil {
		panic(err)
	}

	sort.Sort(model.ByValueDesc(rates.Valutes))

	if err := converter.WriteToJSON(rates.Valutes, config.OutputFile); err != nil {
		panic(err)
	}
}
