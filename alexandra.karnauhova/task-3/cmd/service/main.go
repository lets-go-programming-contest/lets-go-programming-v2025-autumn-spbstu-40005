package main

import (
	"flag"

	"alexandra.karnauhova/task-1/internal/config"
	"alexandra.karnauhova/task-1/internal/parserxml"
)

func main() {
	thisConfig := flag.String("config", "non", "Select a configuration file")
	flag.Parse()

	if *thisConfig == "non" {
		panic("Config file is uncorrect")
	}

	files, err := config.LoadConfig(*thisConfig)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	currencies, err := parserxml.ParseXML(files.InputFile)
	if err != nil {
		panic(err)
	}

}
