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
		return nil, fmt.Errorf("xmlparser parse_file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var target model.CurrencyData
	if err = decoder.Decode(&target); err != nil {
		return nil, fmt.Errorf("xmlparser parse_file %s: %w", path, err)
	}

	return &target, nil
}
