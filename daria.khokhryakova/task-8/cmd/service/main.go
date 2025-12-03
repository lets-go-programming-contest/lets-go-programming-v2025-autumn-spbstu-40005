package main

import (
	"fmt"

	"github.com/DariaKhokhryakova/task-8/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config: %w", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
