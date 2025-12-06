package models

import "sort"

type CurrencyResult struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Currencies []CurrencyResult `xml:"Valute"`
}

func SortByValueDesc(valCurs ValCurs) ([]CurrencyResult, error) {
	results := valCurs.Currencies

	sort.Slice(results, func(i, j int) bool {
		return results[i].Value > results[j].Value
	})

	return results, nil
}
