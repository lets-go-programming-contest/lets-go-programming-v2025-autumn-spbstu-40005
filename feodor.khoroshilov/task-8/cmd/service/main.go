package main

import (
	"fmt"
	"log"

	"feodor.khoroshilov/task-8/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
