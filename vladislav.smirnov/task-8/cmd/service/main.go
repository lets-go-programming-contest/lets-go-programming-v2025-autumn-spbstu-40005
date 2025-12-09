package main

import (
	"fmt"

	"github.com/smirnov-vladislav/task-8/config"
)

func main() {
	config, err := config.ParseConfig()

	if err != nil {
		fmt.Println("error: ", err)

		return
	}

	fmt.Println(config.Environment, " ", config.LogLevel)
}
