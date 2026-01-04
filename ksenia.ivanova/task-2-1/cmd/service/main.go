package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MinAllowed        = 15
	MaxAllowed        = 30
	MinPartsCount     = 2
	ExpectedPartsSize = 2
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		min: MinAllowed,
		max: MaxAllowed,
	}
}

func (tr *TemperatureRange) ApplyConstraint(operator string, value int) {
	switch operator {
	case ">=":
		if value > tr.min {
			tr.min = value
		}
	case "<=":
		if value < tr.max {
			tr.max = value
		}
	}
}

func (tr *TemperatureRange) OptimalTemperature() int {
	if tr.min > tr.max {
		return -1
	}

	return tr.min
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}

	headerParts := strings.Fields(scanner.Text())
	if len(headerParts) < ExpectedPartsSize {
		return
	}

	departmentsCount, _ := strconv.Atoi(headerParts[0])
	employeesPerDept, _ := strconv.Atoi(headerParts[1])

	for range departmentsCount {
		tempRange := NewTemperatureRange()

		for range employeesPerDept {
			if !scanner.Scan() {
				return
			}

			requestParts := strings.Fields(scanner.Text())
			if len(requestParts) < MinPartsCount {
				continue
			}

			operator := requestParts[0]
			value, _ := strconv.Atoi(requestParts[1])

			tempRange.ApplyConstraint(operator, value)
			fmt.Println(tempRange.OptimalTemperature())
		}
	}
}
