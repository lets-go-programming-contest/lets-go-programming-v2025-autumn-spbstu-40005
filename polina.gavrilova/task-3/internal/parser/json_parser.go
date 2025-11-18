package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteJSONData(path string, data interface{}, dirPerm, filePerm os.FileMode) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, jsonData, filePerm); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
