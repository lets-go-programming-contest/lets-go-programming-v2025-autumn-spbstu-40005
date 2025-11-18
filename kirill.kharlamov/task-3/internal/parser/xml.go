package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func LoadXMLData[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("failed to close XML file")
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var data T
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode XML data: %w", err)
	}

	return &data, nil
}
