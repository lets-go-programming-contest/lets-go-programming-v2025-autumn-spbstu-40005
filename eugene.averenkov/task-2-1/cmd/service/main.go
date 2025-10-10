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
	decimalBase    = 10
	minPartsLength = 2
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
	chars := []rune(s)

	for range len(chars) {
		ch := chars[i]
		if ch >= '0' && ch <= '9' {
			num = num*decimalBase + int(ch-'0')
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

	for range departments {
		manager := NewTemperatureManager()

		for range workers {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)

			if len(parts) < minPartsLength {
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
