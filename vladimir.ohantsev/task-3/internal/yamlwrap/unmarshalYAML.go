package yamlwrap

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", op, err))
	}
}

func mustClose(path string, closer io.Closer) {
	must(fmt.Sprintf("close %q", path), closer.Close())
}

func ParseFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open yaml file: %w", err)
	}

	defer mustClose(path, file)

	res, err := Parse[T](file)
	if err != nil {
		return nil, fmt.Errorf("parse yaml file: %w", err)
	}

	return res, nil
}

func Parse[T any](r io.Reader) (*T, error) {
	config := new(T)
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}

	return config, nil
}
