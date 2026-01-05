package main

import (
	"fmt"
	"log"

	"ksenia.ivanova/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Printf("%s %s", cfg.Env, cfg.Log)
}
