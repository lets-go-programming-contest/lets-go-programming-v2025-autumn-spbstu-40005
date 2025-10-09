package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
	minPartsLength = 2
	decimalBase    = 10
)

var (
	ErrUnexpectedEndOfInput = errors.New("unexpected end of input")
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
		return 0, fmt.Errorf("%w", ErrUnexpectedEndOfInput)
	}

	return parseInt(scanner.Text()), nil
}

func parseInt(s string) int {
	num := 0

	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*decimalBase + int(ch-'0')
		}
	}

	return num
}

func processDepartment(scanner *bufio.Scanner, workers int) {
	manager := NewTemperatureManager()

	for range workers {
		if !scanner.Scan() {
			fmt.Println(-1)
			return
		}

		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) < minPartsLength {
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
	for range departments {
		processDepartment(scanner, workers)
	}
}
