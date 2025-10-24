package converter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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

func ParseXMLFile(inputFile string) ([]model.Valute, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading file %s: %w", inputFile, err)
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
		return nil, fmt.Errorf("failed unmarshaling XML: %w", err)
	}

	return rates.Valutes, nil
}

func SortValutes(valutes []model.Valute) {
	sort.Sort(byValueDesc(valutes))
}

func WriteToJSON(valutes []model.Valute, outputFile string) error {
	jsonData, err := json.MarshalIndent(valutes, jsonPrefix, jsonIndent)
	if err != nil {
		return fmt.Errorf("failed marshaling JSON: %w", err)
	}

	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, dirPermissions); err != nil {
		return fmt.Errorf("failed creating directory %s: %w", outputDir, err)
	}

	if err := os.WriteFile(outputFile, jsonData, filePermissions); err != nil {
		return fmt.Errorf("failed writing file %s: %w", outputFile, err)
	}

	return nil
}
