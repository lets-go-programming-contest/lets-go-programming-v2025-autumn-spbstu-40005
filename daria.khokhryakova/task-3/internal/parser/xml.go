package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXMLData[T any](filePath string) (*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	var result T

	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &result, nil
}
