package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXMLData[T any](filePath string) (T, error) {
	var result T

	data, err := os.ReadFile(filePath)
	if err != nil {
		return result, fmt.Errorf("read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)

	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("decode xml: %w", err)
	}

	return result, nil
}
