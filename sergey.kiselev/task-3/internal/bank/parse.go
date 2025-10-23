package bank

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

func ParseXMLFile(filePath string) (*ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error when opening a file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic("error closing XML file")
		}
	}()

	decoder := xml.NewDecoder(file)

	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("XML decoding error: %w", err)
	}

	return &valCurs, nil
}

func (val *value) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var field string
	if err := decoder.DecodeElement(&field, &start); err != nil {
		return fmt.Errorf("cannot decode value: %w", err)
	}

	result, err := strconv.ParseFloat(strings.ReplaceAll(field, ",", "."), 64)
	if err != nil {
		return fmt.Errorf("convert to float: %w", err)
	}

	*val = value(result)

	return nil
}
