package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveAsJSON[T any](obj T, outputPath string, dirPermissions os.FileMode) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create output file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("output file closure failed")
		}
	}()

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(obj); err != nil {
		return fmt.Errorf("JSON encoding failure: %w", err)
	}

	return nil
}
