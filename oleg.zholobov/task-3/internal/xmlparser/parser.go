package xmlparser

import (
	"encoding/xml"
	"errors"
	"io"
	"os"

	"golang.org/x/net/html/charset"
	"oleg.zholobov/task-3/internal/datamodels"
)

var ErrUnsupportedCharset = errors.New("unsupported charset")

func getCharset(charsetLabel string, input io.Reader) (io.Reader, error) {
	encoding, _ := charset.Lookup(charsetLabel)
	if encoding == nil {
		return nil, ErrUnsupportedCharset
	}

	return encoding.NewDecoder().Reader(input), nil
}

func ParseXML(filepath string) ([]datamodels.Valute, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = getCharset

	var valCurs datamodels.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, err
	}

	return valCurs.Valutes, nil
}
