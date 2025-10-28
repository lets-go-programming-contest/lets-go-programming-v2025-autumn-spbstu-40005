package converter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"eugene.averenkov/task-3/internal/currency"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o644
)

func SortByValueDesc(valutes []currency.Valute) {
	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})
}

func WriteToJSON(valutes []currency.Valute, outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("failed close file")
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(valutes); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
