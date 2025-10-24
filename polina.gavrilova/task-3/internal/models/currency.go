package models

import "encoding/xml"

type XMLValCurs struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Valutes []XMLValute `xml:"Valute"`
}

type XMLValute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

type JSONValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}
