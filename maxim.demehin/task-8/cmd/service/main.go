package main

import (
	"fmt"

	"github.com/TvoyBatyA12343/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("failed config loading: %w", err)

		return
	}

	cfg.PrintConfig()
}
