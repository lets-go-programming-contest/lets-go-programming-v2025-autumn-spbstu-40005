package main

import (
	"fmt"

	"gordey.shapkov/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfigFile()
	if err != nil {
		fmt.Println("Err: ", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
