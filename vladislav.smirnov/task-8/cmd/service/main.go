package main

import (
	"fmt"

	"github.com/smirnov-vladislav/task-8/internal/config"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	fmt.Print(config.Environment, " ", config.LogLevel)
}
