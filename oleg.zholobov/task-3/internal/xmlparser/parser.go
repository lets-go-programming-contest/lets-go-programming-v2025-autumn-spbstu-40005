package xmlparser

import (
	"encoding/xml"
	"io"
	"os"

	"oleg.zholobov/task-3/internal/datamodels"
)

func ParseXML(path string) ([]datamodels.Valute, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var valCurs datamodels.ValCurs
	if err := xml.Unmarshal(data, &valCurs); err != nil {
		return nil, err
	}

	return valCurs.Valutes, nil
}
