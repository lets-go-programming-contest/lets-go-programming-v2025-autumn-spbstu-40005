package main

import (
	"fmt"

	"alexandra.karnauhova/task-8/internal/config"
)

func main() {
	cnfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed config: %w", err)

		return
	}

	cnfg.PrintToConfig()
}
