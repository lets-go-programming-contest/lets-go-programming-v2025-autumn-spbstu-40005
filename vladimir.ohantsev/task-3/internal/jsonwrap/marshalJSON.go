package jsonwrap

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func must(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", op, err))
	}
}

func mustClose(path string, closer io.Closer) {
	must(fmt.Sprintf("close %q", path), closer.Close())
}

func Encode[T any](source T, writer io.Writer) error {
	encoder := json.NewEncoder(writer)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(source); err != nil {
		return fmt.Errorf("marshaling: %w", err)
	}

	return nil
}

func EncodeFile[T any](source T, path string, permissions os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer mustClose(path, file)

	if err := Encode(source, file); err != nil {
		return fmt.Errorf("encoding file: %w", err)
	}

	return nil
}
