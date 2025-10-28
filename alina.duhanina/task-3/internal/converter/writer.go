package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"alina.duhanina/task-3/internal/model"
)

func SaveJSON(outputPath string, currencies []model.CurrencyResult) error {
	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing JSON file")
		}
	}()

	if err := os.WriteFile(outputPath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
