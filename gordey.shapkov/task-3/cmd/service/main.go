package main

import (
	"flag"
	"sort"

	"gordey.shapkov/task-3/internal/jsonparsing"
	"gordey.shapkov/task-3/internal/xmlparsing"
	"gordey.shapkov/task-3/internal/yamlparsing"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg, err := yamlparsing.ParseYAMLFile(*configPath)
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
