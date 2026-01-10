package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ksenia.ivanova/task-3/internal/model"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o644
	jsonPrefix      = ""
	jsonIndent      = " "
)

func Dump(data *model.CurrencyData, path string) error {
	jsonData, err := json.MarshalIndent(data.Values, jsonPrefix, jsonIndent)
	if err != nil {
		return fmt.Errorf("Dump Json: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("Dump Json: %w", err)
	}

	if err := os.WriteFile(path, jsonData, filePermissions); err != nil {
		return fmt.Errorf("Dump Json: %w", err)
	}

	return nil
}
