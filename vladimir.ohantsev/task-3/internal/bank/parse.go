package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/P3rCh1/task-3/pkg/must"
	"golang.org/x/net/html/charset"
)

type valueType float64

func (v *valueType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var stringField string
	if err := decoder.DecodeElement(&stringField, &start); err != nil {
		return fmt.Errorf("decode value string: %w", err)
	}

	got, err := strconv.ParseFloat(strings.Replace(stringField, ",", ".", 1), 64)
	if err != nil {
		return fmt.Errorf("invalid value type: %w", err)
	}

	*v = valueType(got)

	return nil
}

func charsetReader(charsetStr string, input io.Reader) (io.Reader, error) {
	return charset.NewReaderLabel(charsetStr, input) //nolint:wrapcheck
}

func ParseXML(r io.Reader) (*Bank, error) {
	decoder := xml.NewDecoder(r)

	decoder.CharsetReader = charsetReader

	bank := new(Bank)
	if err := decoder.Decode(&bank); err != nil {
		return nil, fmt.Errorf("decoding currency bank: %w", err)
	}

	return bank, nil
}

func ParseFileXML(path string) (*Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer must.Close(path, file)

	return ParseXML(file)
}
