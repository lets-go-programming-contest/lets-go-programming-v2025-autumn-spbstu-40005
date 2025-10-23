package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/config"
	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-3/internal/model"
	"golang.org/x/net/html/charset"
)

const (
	jsonPrefix      = ""
	jsonIndent      = "\t"
	dirPermissions  = 0o755
	filePermissions = 0o644
)

type byValueDesc []model.Valute

func (v byValueDesc) Len() int {
	return len(v)
}

func (v byValueDesc) Swap(first, second int) {
	if (first < 0) || (first >= len(v)) {
		panic("first index out of range")
	} else if (second < 0) || (second >= len(v)) {
		panic("second index out of range")
	}

	v[first], v[second] = v[second], v[first]
}

func (v byValueDesc) Less(first, second int) bool {
	if (first < 0) || (first >= len(v)) {
		panic("first index out of range")
	} else if (second < 0) || (second >= len(v)) {
		panic("second index out of range")
	}

	return v[first].Value > v[second].Value
}

func Run(configPath string) error {
	config, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed loading config: %w", err)
	}

	file, err := os.Open(config.InputFile)
	if err != nil {
		return fmt.Errorf("failed reading file %s: %w", config.InputFile, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("failed closing file")
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var rates model.ValCurs
	if err := decoder.Decode(&rates); err != nil {
		return fmt.Errorf("failed unmarshaling XML: %w", err)
	}

	sort.Sort(byValueDesc(rates.Valutes))

	jsonData, err := json.MarshalIndent(rates.Valutes, jsonPrefix, jsonIndent)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %w", err)
	}

	outputDir := filepath.Dir(config.OutputFile)
	if err := os.MkdirAll(outputDir, dirPermissions); err != nil {
		return fmt.Errorf("failed creating directory %s: %w", outputDir, err)
	}

	if err := os.WriteFile(config.OutputFile, jsonData, filePermissions); err != nil {
		return fmt.Errorf("failed writing file %s: %w", config.OutputFile, err)
	}

	return nil
}
