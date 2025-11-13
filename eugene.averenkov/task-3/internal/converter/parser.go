package converter

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ParseXMLFile[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XML file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic("faild close file")
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var result T
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &result, nil
}
