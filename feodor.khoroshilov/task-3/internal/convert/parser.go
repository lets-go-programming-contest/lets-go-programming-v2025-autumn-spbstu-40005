package convert

import (
	"encoding/xml"
	"fmt"
	"os"

	"feodor.khoroshilov/task-3/internal/currency"
	"golang.org/x/net/html/charset"
)

func LoadXMLData(filePath string) ([]currency.Item, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("error closing file")
		}
	}()

	xmlDecoder := xml.NewDecoder(file)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var data currency.ExchangeData
	if err := xmlDecoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding XML: %w", err)
	}

	return data.Items, nil
}
