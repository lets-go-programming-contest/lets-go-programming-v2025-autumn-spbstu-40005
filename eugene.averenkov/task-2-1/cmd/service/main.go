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
	if !tm.IsPossible() {
		return -1
	}
	return tm.minTemp
}

func (tm *TemperatureManager) IsPossible() bool {
	return tm.minTemp <= tm.maxTemp
}

func parseInput(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("unexpected end of input")
	}
	return parseInt(scanner.Text()), nil
}

func parseInt(s string) int {
	num := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		}
	}
	return num
}

func processDepartment(scanner *bufio.Scanner, workers int) {
	manager := NewTemperatureManager()

	for i := 0; i < workers; i++ {
		if !scanner.Scan() {
			fmt.Println(-1)
			return
		}

		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Println(-1)
			continue
		}

		operator := parts[0]
		temp := parseInt(parts[1])

		manager.Update(operator, temp)
		fmt.Println(manager.GetOptimalTemp())
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	departments, err := parseInput(scanner)
	if err != nil {
		fmt.Println(-1)
		return
	}
	workers, err := parseInput(scanner)
	if err != nil {
		fmt.Println(-1)
		return
	}
	for i := 0; i < departments; i++ {
		processDepartment(scanner, workers)
	}
}
