package rates

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"evgeniy.kizhin/task-3/internal/model"
	"golang.org/x/net/html/charset"
)

func LoadRates(path string) ([]model.Currency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read xml: %w", err)
	}

	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.CharsetReader = charset.NewReaderLabel

	var root model.ValCurs

	if err := dec.Decode(&root); err != nil {
		return nil, fmt.Errorf("xml decode: %w", err)
	}

	return root.Valute, nil
}
