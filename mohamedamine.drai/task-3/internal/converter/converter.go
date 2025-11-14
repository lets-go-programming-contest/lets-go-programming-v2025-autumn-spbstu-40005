package converter

import (
	"slices"

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

func MapAndSort[T any](
	items []T,
	getNumCode func(T) int,
	getCharCode func(T) string,
	getValue func(T) float64,
) []CurrencyOutput {
	out := make([]CurrencyOutput, 0, len(items))

	for _, item := range items {
		out = append(out, CurrencyOutput{
			NumCode:  getNumCode(item),
			CharCode: getCharCode(item),
			Value:    getValue(item),
		})
	}

	slices.SortFunc(out, func(a, b CurrencyOutput) int {
		switch {
		case a.Value > b.Value:
			return -1
		case a.Value < b.Value:
			return 1
		default:
			return 0
		}
	})

	return out
}

func (c *CurrencyConverter) ConvertAndSort(valutes []xmlparser.Valute) []CurrencyOutput {
	return MapAndSort(
		valutes,
		func(v xmlparser.Valute) int { return v.NumCode },
		func(v xmlparser.Valute) string { return v.CharCode },
		func(v xmlparser.Valute) float64 { return float64(v.Value) },
	)
}
