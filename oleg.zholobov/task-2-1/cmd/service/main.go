package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation   = errors.New("invalid operator")
	errInvalidTemperature = errors.New("no valid temperature")
	errEmployeeRequest    = errors.New("invalid employee request")
)

const (
	defaultMinTemperature = 15
	defaultMaxTemperature = 30
)

type TemperatureRange struct {
	Min int
	Max int
}

func NewTemperatureRange() TemperatureRange {
	return TemperatureRange{
		Min: defaultMinTemperature,
		Max: defaultMaxTemperature,
	}
}

func (tr *TemperatureRange) Update(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tr.Min {
			tr.Min = temp
		}
	case "<=":
		if temp < tr.Max {
			tr.Max = temp
		}
	default:
		return errInvalidOperation
	}

	return nil
}

func (tr *TemperatureRange) IsValid() bool {
	return tr.Min <= tr.Max
}

func (tr *TemperatureRange) GetOptimalTemperature() (int, error) {
	if !tr.IsValid() {
		return -1, errInvalidTemperature
	}

	return tr.Min, nil
}

func readInput() (int, error) {
	var value int

	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("failed to read input: %w", err)
	}

	return value, nil
}

func readEmployeeRequest() (string, int, error) {
	var operator string

	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return "", 0, fmt.Errorf("failed to read employee request: %w", err)
	}

	return operator, temp, nil
}

func processDepartment(employeeCount int) error {
	temperatureRange := NewTemperatureRange()

	for range employeeCount {
		operator, temp, err := readEmployeeRequest()
		if err != nil {
			return errEmployeeRequest
		}

		err = temperatureRange.Update(operator, temp)
		if err != nil {
			return err
		}

		temp, err = temperatureRange.GetOptimalTemperature()
		if err != nil {
			return err
		}

		fmt.Println(temp)
	}

	return nil
}

func main() {
	departmentCount, err := readInput()
	if err != nil {
		fmt.Println(-1)

		return
	}

	for range departmentCount {
		employeeCount, err := readInput()
		if err != nil {
			fmt.Println(-1)

			return
		}

		if err = processDepartment(employeeCount); err != nil {
			fmt.Println(-1)
		}
	}
}
