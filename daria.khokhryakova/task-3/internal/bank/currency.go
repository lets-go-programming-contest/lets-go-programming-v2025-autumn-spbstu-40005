package bank

type CurrencyResult struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

type ValCurs struct {
	Currencies []CurrencyResult `xml:"Valute"`
}
