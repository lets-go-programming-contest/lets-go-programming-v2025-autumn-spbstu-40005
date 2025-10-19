package bank

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ProcessCurrencies(valCurs *ValCurs) ([]CurrencyResult, error) {
	results := valCurs.Currencies

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func SaveResults(results []CurrencyResult, outputPath string) error {
	dir := filepath.Dir(outputPath)

	err := os.MkdirAll(dir, DirPerm)
	if err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer FileClose(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(results)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
