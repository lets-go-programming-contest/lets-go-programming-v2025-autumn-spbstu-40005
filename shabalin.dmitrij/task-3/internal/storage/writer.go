package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmitei/task-3/internal/models"
)

const (
	DirectoryPermissions = 0o755
	FilePermissions      = 0o600
)

func SaveToJSONFile(destinationPath string, processedCurrencies []models.CurrencyInfo) error {
	destinationDirectory := filepath.Dir(destinationPath)

	if createDirError := os.MkdirAll(destinationDirectory, DirectoryPermissions); createDirError != nil {
		return fmt.Errorf("cannot create output directory %q: %w", destinationDirectory, createDirError)
	}

	jsonEncodedData, encodingError := json.MarshalIndent(processedCurrencies, "", "  ")
	if encodingError != nil {
		return fmt.Errorf("cannot marshal currencies to JSON: %w", encodingError)
	}

	if writeError := os.WriteFile(destinationPath, jsonEncodedData, FilePermissions); writeError != nil {
		return fmt.Errorf("cannot write JSON to file %q: %w", destinationPath, writeError)
	}

	return nil
}
