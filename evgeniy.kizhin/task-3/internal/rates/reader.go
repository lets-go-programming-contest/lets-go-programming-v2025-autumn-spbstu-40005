package rates

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	"evgeniy.kizhin/task-3/internal/model"
	"golang.org/x/net/html/charset"
)

func LoadXML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read xml: %w", err)
	}

	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.CharsetReader = charset.NewReaderLabel

	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("xml decode: %w", err)
	}

	return nil
}

func LoadRates(path string) ([]model.Currency, error) {
	var root model.ValCurs
	if err := LoadXML(path, &root); err != nil {
		return nil, err
	}

	return root.Valute, nil
}
