package parser

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func LoadXMLData[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("file close operation failed")
		}
	}()

	xmlDecoder := xml.NewDecoder(file)
	xmlDecoder.CharsetReader = charset.NewReaderLabel

	var data T
	if err := xmlDecoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("XML parse error: %w", err)
	}

	return &data, nil
}
