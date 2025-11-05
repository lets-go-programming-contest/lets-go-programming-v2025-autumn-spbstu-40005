package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"oleg.zholobov/task-3/internal/datamodels"
)

const dirPermission = 0755

func SaveJSON(path string, currencies []datamodels.Valute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(currencies); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}
