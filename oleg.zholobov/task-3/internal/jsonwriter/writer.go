package jsonwriter

import (
	"encoding/json"
	"os"
	"path/filepath"

	"oleg.zholobov/task-3/internal/datamodels"
)

func SaveJSON(path string, currencies []datamodels.Currency) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(currencies)
}
