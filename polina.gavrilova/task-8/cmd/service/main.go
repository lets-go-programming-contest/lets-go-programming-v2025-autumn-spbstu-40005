package main

import (
	"fmt"

	"polina.gavrilova/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfigFile()
	if err != nil {
		fmt.Println("Err with config loading: %w", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
