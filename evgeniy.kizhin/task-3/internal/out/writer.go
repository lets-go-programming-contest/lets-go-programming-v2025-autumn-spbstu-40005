package out

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func SaveAsJSON(vals any, outPath string) error {
	data, err := json.MarshalIndent(vals, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	dir := filepath.Dir(outPath)

	if dir != "." {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}

	if err := os.WriteFile(outPath, data, filePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
