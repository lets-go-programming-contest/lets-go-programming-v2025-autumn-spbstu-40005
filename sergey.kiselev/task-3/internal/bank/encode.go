package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const permissions = 0o755

func encodeJSON(currencies []Valute, writer io.Writer) error {
	encoder := json.NewEncoder(writer)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(currencies); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}

func EncodeFile(currencies []Valute, outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing output file")
		}
	}()

	return encodeJSON(currencies, file)
}
