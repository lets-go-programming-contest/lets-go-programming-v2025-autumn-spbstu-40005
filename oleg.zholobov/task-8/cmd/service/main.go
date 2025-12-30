package main

import (
	"fmt"

	"oleg.zholobov/task-8/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)

		return
	}

	fmt.Print(cfg.String())
}
