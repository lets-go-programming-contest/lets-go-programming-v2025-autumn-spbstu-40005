package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const fileOpenPermission = 0o755

func Write(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), fileOpenPermission); err != nil {
		return fmt.Errorf("fail to make directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("fail to create file: %w", err)
	}

	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			panic(fmt.Errorf("fail to close file: %w", fileErr))
		}
	}()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("fail to encode: %w", err)
	}

	return nil
}

