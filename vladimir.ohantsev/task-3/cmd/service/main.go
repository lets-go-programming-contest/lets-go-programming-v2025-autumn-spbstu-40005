package main

import (
	"flag"
	"sort"

	"github.com/P3rCh1/task-3/internal/bank"
	"github.com/P3rCh1/task-3/internal/config"
	"github.com/P3rCh1/task-3/pkg/must"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := config.ParseFile(*configPath)

	must.Must("parse config", err)

	bank, err := bank.ParseFileXML(config.Input)

	must.Must("parse input-file", err)

	sort.Slice(
		bank.Currencies,
		func(i, j int) bool {
			return bank.Currencies[i].Value > bank.Currencies[j].Value
		},
	)

	must.Must("encode bank", bank.EncodeFileJSON(config.Output))
}
