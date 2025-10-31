package main

import (
	"flag"

	"feodor.khoroshilov/task-3/internal/config"
	"feodor.khoroshilov/task-3/internal/convert"
)

const (
	configFlagName    = "config"
	configFlagDefault = "config.yaml"
	configFlagUsage   = "path to config file"
)

func main() {
	configPath := flag.String(configFlagName, configFlagDefault, configFlagUsage)
	flag.Parse()

	settings, err := config.LoadSettings(*configPath)
	if err != nil {
		panic(err)
	}

	items, err := convert.LoadXMLData(settings.InputFile)
	if err != nil {
		panic(err)
	}

	convert.SortItemsByRate(items)

	if err := convert.SaveItemsAsJSON(items, settings.OutputFile); err != nil {
		panic(err)
	}
}
