package jsonwriter

import (
	"encoding/json"
	"os"
	"path/filepath"

	"oleg.zholobov/task-3/internal/datamodels"
)

const dirPermission = 0755

func SaveJSON(path string, currencies []datamodels.Valute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
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
