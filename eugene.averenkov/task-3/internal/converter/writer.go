package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteToJSON[T any](data T, outputFile string, dirPermissions, fliePermissions os.FileMode) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("failed close file")
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
