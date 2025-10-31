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

	config_inited, err := config.ConfigInit(*configPath)
	if err != nil {
		panic(err)
	}

	valutes, err := valute.ValCursFromXML(config_inited.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Sort(*valutes)

	if err := jsonstorage.SaveCurrenciesToJSON(*valutes, config_inited.OutputFile); err != nil {
		panic(err)
	}
}
