package jsonparsing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gordey.shapkov/task-3/internal/config"
)

const permissions = 0o755

func ConvertToJSON(valutes []config.Valute) ([]config.Currency, error) {
	currencies := make([]config.Currency, len(valutes))

	for idx, valute := range valutes {
		value, err := config.ConvertFloat(valute.Value)
		if err != nil {
			return nil, fmt.Errorf("convertation failed: %w", err)
		}

		currencies[idx] = config.Currency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
	}

	return currencies, nil
}

func SaveToJSON(currencies []config.Currency, outputPath string) error {
	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal to JSON file: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	if err = os.WriteFile(outputPath, jsonData, permissions); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}
