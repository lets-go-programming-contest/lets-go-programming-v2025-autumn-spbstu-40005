package bank

type Bank struct {
	Currencies []Currency `xml:"Valute"`
}

type Currency struct {
	NumCode  int       `xml:"NumCode"`
	CharCode string    `xml:"CharCode"`
	Value    valueType `xml:"Value"`
}
