package main

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	MinTemperature = 15
	MaxTemperature = 30
)

var (
	ErrInvalidOperator    = errors.New("invalid operator")
	ErrInvalidTemperature = errors.New("invalid temperature")
	ErrNoValidTemperature = errors.New("no valid temperature")
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		min: MinTemperature,
		max: MaxTemperature,
	}
}

func (tr *TemperatureRange) Update(operation string, temperature int) error {
	switch operation {
	case ">=":
		tr.min = max(tr.min, temperature)
	case "<=":
		tr.max = min(tr.max, temperature)
	default:
		return ErrInvalidOperator
	}

	if tr.min > tr.max {
		return ErrNoValidTemperature
	}
	return nil
}

func (tr *TemperatureRange) GetOptimalTemperature() int {
	if tr.min > tr.max {
		return -1
	}
	return tr.min
}

func ParseConstraint(oper string, tempStr string) (string, int, error) {
	if oper != "<=" && oper != ">=" {
		return "", 0, ErrInvalidOperator
	}

	tempInt, err := strconv.Atoi(tempStr)
	if err != nil {
		return "", 0, ErrInvalidTemperature
	}

	if tempInt > MaxTemperature || tempInt < MinTemperature {
		return "", 0, ErrInvalidTemperature
	}

	return oper, tempInt, nil
}

func processDepartment(employeeCount int) []int {
	tr := NewTemperatureRange()
	results := make([]int, 0, employeeCount)

	for i := 0; i < employeeCount; i++ {
		var oper, tempStr string
		if _, err := fmt.Scan(&oper, &tempStr); err != nil {
			results = append(results, -1)
			continue
		}

		operation, temperature, err := ParseConstraint(oper, tempStr)
		if err != nil {
			results = append(results, -1)
			continue
		}

		if err := tr.Update(operation, temperature); err != nil {
			results = append(results, -1)
			continue
		}

		results = append(results, tr.GetOptimalTemperature())
	}

	return results
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println(-1)
		return
	}

	for i := 0; i < departments; i++ {
		var employees int

		_, err := fmt.Scan(&employees)
		if err != nil {
			fmt.Println(-1)
			return
		}

		results := processDepartment(employees)
		for _, result := range results {
			fmt.Println(result)
		}
	}
}
