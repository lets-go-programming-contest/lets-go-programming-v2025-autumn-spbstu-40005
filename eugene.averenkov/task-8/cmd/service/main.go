package main

import (
	"fmt"

	"eugene.averenkov/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v", err)

		return
	}

	fmt.Print(cfg.Env, " ", cfg.LogLevel)
}
