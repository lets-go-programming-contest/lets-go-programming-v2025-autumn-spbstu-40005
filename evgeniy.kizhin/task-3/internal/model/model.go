package model

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type ValCurs struct {
	XMLName xml.Name   `xml:"ValCurs"`
	Date    string     `xml:"Date,attr"`
	Name    string     `xml:"name,attr"`
	Valute  []Currency `xml:"Valute"`
}

type Amount float64

func (a *Amount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	s = strings.ReplaceAll(s, ",", ".")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*a = Amount(v)
	return nil
}

type Currency struct {
	NumCode  int    `xml:"NumCode" json:"num_code"`
	CharCode string `xml:"CharCode" json:"char_code"`
	Value    Amount `xml:"Value" json:"value"`
}
