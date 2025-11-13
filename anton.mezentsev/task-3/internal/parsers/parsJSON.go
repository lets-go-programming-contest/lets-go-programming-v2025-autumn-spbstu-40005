package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveToJSON(data any, filePath string, directoryPermissions os.FileMode) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, directoryPermissions); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("closing file: " + err.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}

	return nil
}
