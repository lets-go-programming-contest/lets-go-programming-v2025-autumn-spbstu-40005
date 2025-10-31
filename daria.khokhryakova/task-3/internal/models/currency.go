package models

type CurrencyResult struct {
	NumCode  int           `xml:"NumCode"  json:"num_code"`
	CharCode string        `xml:"CharCode" json:"char_code"`
	Value    CurrencyValue `xml:"Value"    json:"value"`
}

type ValCurs struct {
	Currencies []CurrencyResult `xml:"Valute"`
}
