package main

import (
	"flag"
	"sort"

	"evgeniy.kizhin/task-3/internal/cfg"
	"evgeniy.kizhin/task-3/internal/out"
	"evgeniy.kizhin/task-3/internal/rates"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to yaml configuration")
	flag.Parse()

	conf, err := cfg.Load(*cfgPath)
	if err != nil {
		panic(err)
	}

	rates, err := rates.LoadRates(conf.InputFile)
	if err != nil {
		panic(err)
	}

	sort.Slice(rates, func(i, j int) bool {
		return float64(rates[i].Value) > float64(rates[j].Value)
	})

	if err := out.SaveAsJSON(rates, conf.OutputFile); err != nil {
		panic(err)
	}
}
