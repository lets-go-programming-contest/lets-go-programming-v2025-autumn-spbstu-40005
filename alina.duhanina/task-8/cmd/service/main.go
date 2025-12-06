package main

import (
	"fmt"
	"alina.duhanina/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("failed loading: %w", err)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
