package main

import (
	"alina.duhanina/task-8/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("failed config loading: %w", err)

		return
	}

	cfg.PrintConfig()
}}
