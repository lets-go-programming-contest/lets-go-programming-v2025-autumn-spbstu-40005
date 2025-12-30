package writer

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
)

func CreateDirectory(directory string, directoryPermissions fs.FileMode) error {
	if err := os.MkdirAll(directory, directoryPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

func SaveToJSON(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
