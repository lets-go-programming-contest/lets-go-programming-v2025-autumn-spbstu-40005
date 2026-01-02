package main

import (
	"fmt"

	"sergey.dribas/task-8/internal/config"
)

func main() {
	if cnfg, err := config.GetConfig(); err != nil {
		fmt.Println("failed config: %w", err)
		return
	} else {
		cnfg.PrintToConfig()
	}
}
