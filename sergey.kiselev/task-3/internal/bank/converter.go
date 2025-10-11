package bank

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ConvertAndSort(valutes []Valute) ([]Currency, error) {
	currencies := make([]Currency, 0, len(valutes))

	for _, valute := range valutes {
		value, err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."), 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value for %s: %w", valute.CharCode, err)
		}

		currencies = append(currencies, Currency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies, nil
}
