package jsonstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sergey.dribas/task-3/internal/data"
)

const (
	defaultFilePerm = 0600
	defaultDirPerm  = 0755
)

type CurrencyJSON struct {
	NumCode  int             `json:"num_code"`
	CharCode string          `json:"char_code"`
	Value    json.RawMessage `json:"value"`
}

func SaveCurrenciesToJSON(currencies valute.ValCurs, filename string) error {
	result := make([]CurrencyJSON, 0, len(currencies.Valutes))
	dir := filepath.Dir(filename)

	for _, currency := range currencies.Valutes {
		var (
			numCode int
			err     error
		)

		if currency.NumCode != "" {
			numCode, err = strconv.Atoi(currency.NumCode)
			if err != nil {
				return fmt.Errorf("error cast to int: %w", err)
			}
		}

		value := strings.Replace(currency.Value, ",", ".", 1)
		rawValue := json.RawMessage(value)
		result = append(result, CurrencyJSON{
			NumCode:  numCode,
			CharCode: currency.CharCode,
			Value:    rawValue,
		})
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	if err := os.MkdirAll(dir, defaultDirPerm); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	if err = os.WriteFile(filename, jsonData, defaultFilePerm); err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	return nil
}
