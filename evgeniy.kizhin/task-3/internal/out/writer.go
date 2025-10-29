package out

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"evgeniy.kizhin/task-3/internal/model"
)

const dirPerm = 0o755
const filePerm = 0o644

func SaveAsJSON(vals []model.Currency, outPath string) error {
	b, err := json.MarshalIndent(vals, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	dir := filepath.Dir(outPath)
	if dir != "." {
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}
	}
	if err := os.WriteFile(outPath, b, filePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
