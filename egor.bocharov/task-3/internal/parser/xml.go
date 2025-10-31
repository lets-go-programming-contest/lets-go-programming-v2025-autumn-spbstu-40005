package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXMLFile[T any](filePath string) (*T, error) {
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

	var result T
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("XML decoding error: %w", err)
	}

	return &result, nil
}
