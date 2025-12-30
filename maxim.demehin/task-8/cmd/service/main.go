package main

import (
	"fmt"

	"github.com/TvoyBatyA12343/task-8/internal/config"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Println("failed config loading: %w", err)

		return
	}

	cfg.PrintConfig()
}
