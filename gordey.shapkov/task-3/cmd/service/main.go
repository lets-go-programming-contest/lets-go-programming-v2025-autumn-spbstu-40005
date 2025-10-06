package main

import (
	"flag"
	"fmt"
	"sort"

	"gordey.shapkov/task-3/internal/parsers"
)

func main() {
	configPath := flag.String("config", "", "YAML file required")
	flag.Parse()

	cfg, err := parsers.ParseConfigFile(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := parsers.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currencies, err := parsers.ConvertToJSON(valCurs.Valutes)
	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	if err = parsers.SaveToJSON(currencies, cfg.OutputFile); err != nil {
		fmt.Println(err)
	}
}
