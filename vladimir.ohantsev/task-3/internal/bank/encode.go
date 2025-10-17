package bank

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/P3rCh1/task-3/pkg/must"
)

const permissions = 0o755

func (b *Bank) EncodeJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)

	encoder.SetIndent("", "  ")

	if err := encoder.Encode(b.Currencies); err != nil {
		return fmt.Errorf("encoding bank: %w", err)
	}

	return nil
}

func (b *Bank) EncodeFileJSON(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, permissions); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer must.Close(path, file)

	if err := b.EncodeJSON(file); err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	return nil
}
