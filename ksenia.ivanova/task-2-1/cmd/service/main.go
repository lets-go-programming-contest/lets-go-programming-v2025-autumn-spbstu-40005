package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	minLimit = 15
	maxLimit = 30
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
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	departmentsCount, _ := strconv.Atoi(scanner.Text())

	if !scanner.Scan() {
		return
	}
	employeesPerDept, _ := strconv.Atoi(scanner.Text())

	for range departmentsCount {
		tempRange := NewTemperatureRange(minLimit, maxLimit)

		for range employeesPerDept {
			if !scanner.Scan() {
				return
			}
			operator := scanner.Text()

			if !scanner.Scan() {
				return
			}
			tempValue, _ := strconv.Atoi(scanner.Text())

			UpdateTemperature(operator, tempValue, tempRange)
			fmt.Println(tempRange.GetResult())
		}
	}
}
