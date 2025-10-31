package main

import (
	"flag"
	"path/filepath"

	"alexandra.karnauhova/task-1/internal/config"
	"alexandra.karnauhova/task-1/internal/parserxml"
	"alexandra.karnauhova/task-1/internal/sorter"
	"alexandra.karnauhova/task-1/internal/writer"
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

	valArray, err := parserxml.ParseXML(files.InputFile)
	if err != nil {
		panic(err)
	}

	sortedValutes := sorter.SortByValueDesc(valArray.Valutes)

	outputDir := filepath.Dir(files.OutputFile)
	if err := writer.CreateDirectory(outputDir); err != nil {
		panic("Failed to create output directory: " + err.Error())
	}

	if err := writer.ParseJson(sortedValutes, files.OutputFile); err != nil {
		panic("Failed to save JSON: " + err.Error())
	}
}
