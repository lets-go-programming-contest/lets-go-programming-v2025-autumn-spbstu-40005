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

func SaveToJSON(data any, filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingDir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingFile, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(ErrClosingFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("%w: %v", ErrEncodingJSON, err)
	}

	return nil
}
