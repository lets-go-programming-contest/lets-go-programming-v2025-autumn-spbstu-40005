package main

import (
	"flag"
	"sort"

	"sergey.dribas/task-3/internal/config"
	"sergey.dribas/task-3/internal/data"
	"sergey.dribas/task-3/internal/json"
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

	valutes, err := valute.ValCursFromXML(config.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Sort(*valutes)

	if err := jsonstorage.SaveCurrenciesToJSON(*valutes, config.OutputFile); err != nil {
		panic(err)
	}
}
