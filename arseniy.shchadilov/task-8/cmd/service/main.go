package main

import (
	"fmt"

	"github.com/arseniy.shchadilov/task-8/config"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
