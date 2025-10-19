package bank

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ParseXMLData(filePath string) (*ValCurs, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var valCurs ValCurs

	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &valCurs, nil
}
