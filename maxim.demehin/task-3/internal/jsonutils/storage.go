package jsonutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const permission = 0o755

func SaveToFile[T any](data T, output string) error {
	dir := filepath.Dir(output)

	err := os.MkdirAll(dir, permission)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshal to JSON: %w", err)
	}

	err = os.WriteFile(output, jsonData, permission)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
