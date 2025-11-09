package converter

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"mohamedamine.drai/task-3/internal/xmlparser"
)

type CurrencyOutput struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type CurrencyConverter struct{}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{}
}

func (c *CurrencyConverter) ConvertAndSort(currencies []xmlparser.Currency) []CurrencyOutput {
	output := make([]CurrencyOutput, 0, len(currencies))

	for _, currency := range currencies {
		converted, ok := c.convertCurrency(currency)
		if ok {
			output = append(output, converted)
		}
	}

	c.sortByValueDescending(output)

	return output
}

func (c *CurrencyConverter) convertCurrency(currency xmlparser.Currency) (CurrencyOutput, bool) {
	value, err := c.parseValue(currency.Value)
	if err != nil {
		return CurrencyOutput{
			NumCode:  0,
			CharCode: "",
			Value:    0,
		}, false
	}

	numCode, err := c.parseNumCode(currency.NumCode)
	if err != nil {
		return CurrencyOutput{
			NumCode:  0,
			CharCode: "",
			Value:    0,
		}, false
	}

	return CurrencyOutput{
		NumCode:  numCode,
		CharCode: strings.TrimSpace(currency.CharCode),
		Value:    value,
	}, true
}

func (c *CurrencyConverter) parseValue(valueStr string) (float64, error) {
	normalized := strings.ReplaceAll(valueStr, ",", ".")
	value, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0, fmt.Errorf("parse value %q: %w", normalized, err)
	}

	return value, nil
}

func (c *CurrencyConverter) parseNumCode(numCodeStr string) (int, error) {
	numCode, err := strconv.Atoi(strings.TrimSpace(numCodeStr))
	if err != nil {
		return 0, fmt.Errorf("parse num code %q: %w", numCodeStr, err)
	}

	return numCode, nil
}

func (c *CurrencyConverter) sortByValueDescending(currencies []CurrencyOutput) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}
