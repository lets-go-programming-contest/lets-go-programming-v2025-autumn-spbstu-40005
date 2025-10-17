package main

import (
	"flag"
	"fmt"
	"sort"

	"gordey.shapkov/task-3/internal/jsonparsing"
	"gordey.shapkov/task-3/internal/xmlparsing"
	"gordey.shapkov/task-3/internal/yamlparsing"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred", err)

			return
		}
	}()

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

	currencies, err := jsonparsing.ConvertToJSON(valCurs.Valutes)
	if err != nil {
		panic(err)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	if err = jsonparsing.SaveToJSON(currencies, cfg.OutputFile); err != nil {
		panic(err)
	}
}
