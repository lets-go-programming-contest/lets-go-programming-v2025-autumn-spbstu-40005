package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrCreatingDir  = errors.New("error creating directory")
	ErrCreatingFile = errors.New("error creating file")
	ErrEncodingJSON = errors.New("error encoding json")
	ErrClosingFile  = errors.New("error closing file")
)

func SaveToJSON(data any, filePath string, directoryPermissions os.FileMode) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, directoryPermissions); err != nil {
		return fmt.Errorf("%w: %w", ErrCreatingDir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreatingFile, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(ErrClosingFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("%w: %w", ErrEncodingJSON, err)
	}

	return nil
}
