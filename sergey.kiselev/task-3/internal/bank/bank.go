package bank

type ValCurs struct {
	Date   string   `xml:"Date,attr"`
	Name   string   `xml:"name,attr"`
	Valute []Valute `xml:"Valute"`
}

type value float64

type Valute struct {
	NumCode  int    `json:"num_code"  xml:"NumCode"`
	CharCode string `json:"char_code" xml:"CharCode"`
	Value    value  `json:"value"     xml:"Value"`
}
