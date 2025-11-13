package main

import (
	"flag"
	"sort"

	"gordey.shapkov/task-3/internal/config"
	"gordey.shapkov/task-3/internal/jsonparsing"
	"gordey.shapkov/task-3/internal/xmlparsing"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.ParseConfigFile(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := xmlparsing.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(valCurs.Valutes, func(i, j int) bool {
		return valCurs.Valutes[i].Value > valCurs.Valutes[j].Value
	})

	if err = jsonparsing.SaveToJSON(valCurs.Valutes, cfg.OutputFile); err != nil {
		panic(err)
	}
}
