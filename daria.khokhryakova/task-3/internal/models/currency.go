package models

type CurrencyResult struct {
	NumCode  int           `json:"num_code"  xml:"NumCode"`
	CharCode string        `json:"char_code" xml:"CharCode"`
	Value    CurrencyValue `json:"value"     xml:"Value"`
}

type ValCurs struct {
	Currencies []CurrencyResult `xml:"Valute"`
}
