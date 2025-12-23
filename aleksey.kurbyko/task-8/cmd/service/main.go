package main

import (
	"fmt"

	"aleksey.kurbyko/task-8/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
