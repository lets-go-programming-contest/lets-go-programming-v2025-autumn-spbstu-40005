package xmlwrap

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

func must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", op, err))
	}
}

func mustClose(path string, closer io.Closer) {
	must(fmt.Sprintf("close %q", path), closer.Close())
}

func Parse[T any](r io.Reader) (*T, error) {
	decoder := xml.NewDecoder(r)

	decoder.CharsetReader = charset.NewReaderLabel

	result := new(T)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}

	return result, nil
}

func ParseFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open input file: %w", err)
	}

	defer mustClose(path, file)

	res, err := Parse[T](file)
	if err != nil {
		return nil, fmt.Errorf("parse xml file: %w", err)
	}

	return res, nil
}
