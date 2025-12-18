package main

import (
	"flag"
	"sort"

	"github.com/Rassoha/lets-go-programming-v2025-autumn-spbstu-40005/sergey.dribas/task-3/internal/config"
	"github.com/Rassoha/lets-go-programming-v2025-autumn-spbstu-40005/sergey.dribas/task-3/internal/jsonstorage"
	"github.com/Rassoha/lets-go-programming-v2025-autumn-spbstu-40005/sergey.dribas/task-3/internal/valute"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	config, err := config.ConfigInit(*configPath)
	if err != nil {
		panic(err)
	}

	valutes := &valute.ValCurs{
		Valutes: []valute.Valute{},
	}

	err = valute.ParseXML(config.InputFile, valutes)
	if err != nil {
		panic(err)
	}

	sort.Sort(*valutes)

	if err := jsonstorage.SaveCurrenciesToJSON(*valutes, config.OutputFile); err != nil {
		panic(err)
	}
}
