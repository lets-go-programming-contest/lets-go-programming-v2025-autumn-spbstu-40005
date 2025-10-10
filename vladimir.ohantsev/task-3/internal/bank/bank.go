package bank

type Bank struct {
	Currencies []Currency `json:"valute" xml:"Valute"`
}

type Currency struct {
	NumCode  int       `json:"num_code"  xml:"NumCode"`
	CharCode string    `json:"char_code" xml:"CharCode"`
	Value    valueType `json:"value"     xml:"Value"`
}
