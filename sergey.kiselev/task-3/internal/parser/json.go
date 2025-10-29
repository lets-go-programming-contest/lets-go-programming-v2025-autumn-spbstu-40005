package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const permissions = 0o755

func EncodeFile[T any](data T, outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing output file")
		}
	}()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}
