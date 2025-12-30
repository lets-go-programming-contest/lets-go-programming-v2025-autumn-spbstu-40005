package main

import (
	"fmt"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-8/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("when loading config: %w", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
