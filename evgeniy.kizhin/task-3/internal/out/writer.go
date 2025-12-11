package out

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SaveOpts struct {
	DirPerm  os.FileMode
	FilePerm os.FileMode
}

func SaveAsJSON(opts SaveOpts, vals any, outPath string) error {
	data, err := json.MarshalIndent(vals, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	dir := filepath.Dir(outPath)

	if err := os.MkdirAll(dir, opts.DirPerm); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(outPath, data, opts.FilePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
