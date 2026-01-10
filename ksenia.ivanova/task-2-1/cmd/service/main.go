package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	minLimit      = 15
	maxLimit      = 30
	minPartsCount = 2
	expectedParts = 2
)

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange(minVal, maxVal int) *TemperatureRange {
	return &TemperatureRange{
		minTemp: minVal,
		maxTemp: maxVal,
	}
}

func (temp *TemperatureRange) IsValid() bool {
	return temp.minTemp <= temp.maxTemp
}

func (temp *TemperatureRange) GetResult() int {
	if !temp.IsValid() {
		return -1
	}

	return temp.minTemp
}

func UpdateTemperature(operator string, temperature int, tempRange *TemperatureRange) {
	switch operator {
	case ">=":
		if temperature > tempRange.minTemp {
			tempRange.minTemp = temperature
		}
	case "<=":
		if temperature < tempRange.maxTemp {
			tempRange.maxTemp = temperature
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}

	departmentsCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || departmentsCount <= 0 {
		return
	}

	for range departmentsCount {
		if !scanner.Scan() {
			return
		}

		employeesCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil || employeesCount <= 0 {
			return
		}

		tempRange := NewTemperatureRange(minLimit, maxLimit)

		for range employeesCount {
			if !scanner.Scan() {
				return
			}

			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)

			if len(parts) < minPartsCount {
				continue
			}

			operator := parts[0]
			value, parseErr := strconv.Atoi(parts[1])

			if parseErr != nil {
				continue
			}

			UpdateTemperature(operator, value, tempRange)
			fmt.Println(tempRange.GetResult())
		}
	}
}
