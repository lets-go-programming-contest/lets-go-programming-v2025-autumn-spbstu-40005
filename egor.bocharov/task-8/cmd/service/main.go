package main

import (
	"fmt"

	"egor.bocharov/task-8/config"
)

func main() {
	cfg, err := config.GetConfigFile()
	if err != nil {
		fmt.Println("Err: ", err)

		return
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
