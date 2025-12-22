package main

import (
	"fmt"
	"oleg.zholobov/task-8/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load config: %v", err)

		return
	}
	
	fmt.Println(cfg.String())
}
