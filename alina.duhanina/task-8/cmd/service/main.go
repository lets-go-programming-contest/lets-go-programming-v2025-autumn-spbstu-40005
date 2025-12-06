package main

import (
	"alina.duhanina/task-8/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed loading: %w", err)

		return
	}

	fmt.Print(cfg.Environment + " " + cfg.LogLevel)
}
