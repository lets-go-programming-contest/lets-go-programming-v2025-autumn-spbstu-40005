package main

import (
	"flag"
	"sort"

	"github.com/deonik3/task-3/internal/bank"
	"github.com/deonik3/task-3/internal/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := config.ParseFile(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := bank.ParseXMLFile(config.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(valCurs.Valute, func(i, j int) bool {
		return valCurs.Valute[i].Value > valCurs.Valute[j].Value
	})

	if err = bank.EncodeFile(valCurs.Valute, config.OutputFile); err != nil {
		panic(err)
	}
}
