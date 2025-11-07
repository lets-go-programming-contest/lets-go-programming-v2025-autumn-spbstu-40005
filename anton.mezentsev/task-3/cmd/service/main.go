package main

import (
	"flag"
	"os"
	"sort"

	"anton.mezentsev/task-3/internal/bank"
	"anton.mezentsev/task-3/internal/config"
	"anton.mezentsev/task-3/internal/parsers"
)

func main() {
	configfilePath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg, err := parsers.ParseYAML[config.Config](*configfilePath)
	if err != nil {
		panic(err)
	}

	valCurs, err := parsers.ParseXML[bank.ValCurs](cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(valCurs.Valute, func(index1, index2 int) bool {
		return valCurs.Valute[index1].Value > valCurs.Valute[index2].Value
	})

	dirPerms := os.FileMode(cfg.DirPerms)
	if err = parsers.SaveToJSON(valCurs.Valute, cfg.OutputFile, dirPerms); err != nil {
		panic(err)
	}
}
