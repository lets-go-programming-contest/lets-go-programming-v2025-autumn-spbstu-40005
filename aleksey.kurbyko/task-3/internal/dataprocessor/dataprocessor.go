package dataprocessor

import (
	"strconv"
	"strings"

	"github.com/aleksey.kurbyko/task-3/internal/currencyhandler"
)

type FilePaths struct {
	Input  string `yaml:"input-file"`
	Output string `yaml:"output-file"`
}

type JSONCurrency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func processCurrency(currency currencyhandler.CurrencyItem) JSONCurrency {
	var result JSONCurrency
	var parseError error

	if currency.NumCode != "" {
		result.NumCode, parseError = strconv.Atoi(currency.NumCode)
		if parseError != nil {
			panic("Number code conversion error")
		}
	}

	cleanValue := strings.ReplaceAll(currency.Value, ",", ".")
	result.Value, parseError = strconv.ParseFloat(cleanValue, 64)
	if parseError != nil {
		panic("Value conversion error")
	}

	result.CharCode = currency.CharCode
	return result
}

func ConvertToJSON(data currencyhandler.CurrencyList) []JSONCurrency {
	converted := make([]JSONCurrency, 0, len(data.Items))
	for _, item := range data.Items {
		converted = append(converted, processCurrency(item))
	}
	return converted
}
