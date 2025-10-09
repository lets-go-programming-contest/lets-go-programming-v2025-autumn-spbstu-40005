package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
)

type TemperatureManager struct {
	minTemp int
	maxTemp int
}

func NewTemperatureManager() *TemperatureManager {
	return &TemperatureManager{
		minTemp: InitialMinTemp,
		maxTemp: InitialMaxTemp,
	}
}

func (tm *TemperatureManager) Update(operator string, temp int) {
	switch operator {
	case ">=":
		if temp > tm.minTemp {
			tm.minTemp = temp
		}
	case "<=":
		if temp < tm.maxTemp {
			tm.maxTemp = temp
		}
	}
}

func (tm *TemperatureManager) GetOptimalTemp() int {
	if tm.minTemp > tm.maxTemp {
		return -1
	}

	return tm.minTemp
}

func stringToInt(s string) int {
	num := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		}
	}
	return num
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	departments := stringToInt(scanner.Text())

	scanner.Scan()
	workers := stringToInt(scanner.Text())

	for i := 0; i < departments; i++ {
		manager := NewTemperatureManager()

		for j := 0; j < workers; j++ {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)

			if len(parts) < 2 {
				fmt.Println(-1)
				continue
			}

			operator := parts[0]
			temp := stringToInt(parts[1])

			manager.Update(operator, temp)
			fmt.Println(manager.GetOptimalTemp())
		}
	}
}

