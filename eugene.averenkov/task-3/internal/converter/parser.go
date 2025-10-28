package converter

import (
	"currency-converter/internal/currency"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"os"
)

func ParseXMLFile(filePath string) ([]currency.Valute, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close file: %v\n", closeErr)
		}
	}()
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs currency.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return valCurs.Valutes, nil
}
