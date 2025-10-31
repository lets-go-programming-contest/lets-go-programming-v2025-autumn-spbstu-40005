package models

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CurrencyValue float64

func (cv *CurrencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var valueStr string
	if err := d.DecodeElement(&valueStr, &start); err != nil {
		return err
	}

	value, err := convertValue(valueStr)
	if err != nil {
		return err
	}

	*cv = CurrencyValue(value)

	return nil
}

func convertValue(valueStr string) (float64, error) {
	normalizStr := strings.ReplaceAll(valueStr, ",", ".")

	value, err := strconv.ParseFloat(normalizStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}
	return value, nil
}
