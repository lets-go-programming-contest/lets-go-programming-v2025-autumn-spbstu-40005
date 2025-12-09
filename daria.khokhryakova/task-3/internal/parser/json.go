package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SaveJSONResults[T any](results []T, outputPath string, dirPerm os.FileMode) error {
	dir := filepath.Dir(outputPath)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer func() {
		panicErr(file.Close())
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(results)
	if err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}
