package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/DariaKhokhryakova/task-3/internal/models"
)

const dirPerm = 0o755

func ProcessCurrencies(valCurs *models.ValCurs) ([]models.CurrencyResult, error) {
	results := valCurs.Currencies

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func FileClose(file *os.File) {
	if file != nil {
		err := file.Close()
		panicErr(err)
	}
}

func SaveJSONResults(results []models.CurrencyResult, outputPath string) error {
	dir := filepath.Dir(outputPath)

	err := os.MkdirAll(dir, dirPerm)
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
