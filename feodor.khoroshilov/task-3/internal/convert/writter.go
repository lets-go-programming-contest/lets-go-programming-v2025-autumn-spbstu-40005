// writter.go
package convert

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"feodor.khoroshilov/task-3/internal/currency"
)

const (
	dirPerms  = 0o755
	filePerms = 0o644
)

func SortItemsByRate(items []currency.Item) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].RateValue > items[j].RateValue
	})
}

func SaveItemsAsJSON(items []currency.Item, outputPath string) error {
	dirName := filepath.Dir(outputPath)
	if err := os.MkdirAll(dirName, dirPerms); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if closeErr := output.Close(); closeErr != nil {
			panic("error closing file")
		}
	}()

	jsonEncoder := json.NewEncoder(output)
	jsonEncoder.SetIndent("", "  ")

	if err := jsonEncoder.Encode(items); err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}