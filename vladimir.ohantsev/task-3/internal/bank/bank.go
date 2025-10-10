package bank

type Bank struct {
	Currencies []Currency `xml:"Valute" json:"valute"`
}

type Currency struct {
	NumCode  int       `xml:"NumCode" json:"num_code"`
	CharCode string    `xml:"CharCode" json:"char_code"`
	Value    valueType `xml:"Value" json:"value"`
}
