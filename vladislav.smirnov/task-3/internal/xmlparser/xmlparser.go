package xmlparser

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

func ReadXML(path string, target any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail to open file %q: %w", path, err)
	}

	defer func() {
		if fileErr := file.Close(); fileErr != nil {
			panic(fmt.Errorf("fail to close file %q: %w", path, fileErr))
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("fail to decode %q: %w", path, err)
	}

	return nil
}
