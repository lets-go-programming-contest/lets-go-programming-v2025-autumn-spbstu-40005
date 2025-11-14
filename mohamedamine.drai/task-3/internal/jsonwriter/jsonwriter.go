package jsonwriter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultDirPerm  = 0o755
	DefaultFilePerm = 0o600
)

func SaveToJSON(data any, path string) error {
	return SaveToJSONWithPerms(data, path, DefaultDirPerm, DefaultFilePerm)
}

func SaveToJSONWithPerms(data any, path string, dirPerm os.FileMode, filePerm os.FileMode) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	bytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(path, bytes, filePerm); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}
