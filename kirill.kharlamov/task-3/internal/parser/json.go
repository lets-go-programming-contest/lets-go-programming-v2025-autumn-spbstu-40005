package parser

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const defaultDirPermissions = 0o755

func SaveAsJSON[T any](data T, outputPath string, dirPermissions *uint32) error {
	dir := filepath.Dir(outputPath)

	var perms fs.FileMode = defaultDirPermissions
	if dirPermissions != nil {
		perms = fs.FileMode(*dirPermissions)
	}

	if err := os.MkdirAll(dir, perms); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("failed to close output file")
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON data: %w", err)
	}

	return nil
}
