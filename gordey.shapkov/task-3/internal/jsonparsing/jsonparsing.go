package jsonparsing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gordey.shapkov/task-3/internal/models"
)

const permissions = 0o755

func SaveToJSON(currencies []models.Currency, outputPath string) error {
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
