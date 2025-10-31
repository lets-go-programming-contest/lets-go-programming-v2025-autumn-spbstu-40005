package writer

import (
	"encoding/json"
	"fmt"
	"os"

	"alexandra.karnauhova/task-3/internal/data"
)

const directoryPermissions = 0o755

func CreateDirectory(directory string) error {
	return os.MkdirAll(directory, directoryPermissions)
}

func SaveToJSON(valutes []data.Valute, filePath string) error {
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

	if err := encoder.Encode(valutes); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
