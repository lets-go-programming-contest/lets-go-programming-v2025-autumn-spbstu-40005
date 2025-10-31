package jsonstorage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"sergey.dribas/task-3/internal/data"
)

type CurrencyJSON struct {
	NumCode  int             `json:"num_code"`
	CharCode string          `json:"char_code"`
	Value    json.RawMessage `json:"value"`
}

func SaveCurrenciesToJSON(currencies valute.ValCurs, filename string) error {
	var result []CurrencyJSON
	dir := filepath.Dir(filename)
	for _, currency := range currencies.Valutes {
		if numCode, err := strconv.Atoi(currency.NumCode); err != nil {
			return err
		} else {
			value := strings.Replace(currency.Value, ",", ".", 1)
			rawValue := json.RawMessage(value)
			result = append(result, CurrencyJSON{
				NumCode:  numCode,
				CharCode: currency.CharCode,
				Value:    rawValue,
			})
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
