package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"

	"ksenia.ivanova/task-3/internal/model"
)

func ParseFile(path string) (*model.CurrencyData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file %s: %v\n", path, closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var target model.CurrencyData
	if err = decoder.Decode(&target); err != nil {
		return nil, fmt.Errorf("parse file %s: %w", path, err)
	}

	return &target, nil
}
