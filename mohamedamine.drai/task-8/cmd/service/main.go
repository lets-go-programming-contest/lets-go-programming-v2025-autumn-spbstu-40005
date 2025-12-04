package main

import (
	"fmt"
	"github.com/aminedraii19/task-8/internal/config"
)

func main() {
	fmt.Println(config.Cfg.Environment, config.Cfg.LogLevel)
}
