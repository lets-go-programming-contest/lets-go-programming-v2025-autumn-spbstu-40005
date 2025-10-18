package jsonutils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TvoyBatyA12343/task-3/internal/bank"
)

const permission = 0o755

func SaveValutesToFile(valutes []bank.Valute, output string) error {
	dir := filepath.Dir(output)

	err := os.MkdirAll(dir, permission)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(valutes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshal to JSON: %w", err)
	}

	err = os.WriteFile(output, jsonData, permission)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
