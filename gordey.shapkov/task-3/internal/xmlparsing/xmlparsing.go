package xmlparsing

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/charmap"
	"gordey.shapkov/task-3/internal/config"
)

func ParseXMLFile(path string) (*config.ValCurs, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = createCharsetReader

	valCurs := &config.ValCurs{Date: "", Name: "", Valutes: nil}
	if err = decoder.Decode(valCurs); err != nil {
		return nil, fmt.Errorf("cannot decode file: %w", err)
	}

	return valCurs, nil
}

func createCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1251" {
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	return input, nil
}
