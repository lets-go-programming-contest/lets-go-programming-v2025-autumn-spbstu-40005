package main

import (
	"errors"
	"fmt"
)

var (
	errInvalidOperation = errors.New("invalid operator")
	errEmployeeRequest  = errors.New("invalid employee request")
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

func (tr *TemperatureRange) GetOptimalTemperature() int {
	if !tr.IsValid() {
		return -1
	}

	return tr.Min
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
	var (
		operator string
		temp     int
	)

	if _, err := fmt.Scan(&operator, &temp); err != nil {
		return "", 0, fmt.Errorf("failed to read employee request: %w", err)
	}

	return operator, temp, nil
}

func processDepartment() error {
	employeeCount, err := readInput()
	if err != nil {
		return fmt.Errorf("error reading number of employee: %w", err)
	}

	temperatureRange := NewTemperatureRange()

	for range employeeCount {
		operator, temp, err := readEmployeeRequest()
		if err != nil {
			return errEmployeeRequest
		}

		if err = temperatureRange.Update(operator, temp); err != nil {
			return err
		}

		temp = temperatureRange.GetOptimalTemperature()
		fmt.Println(temp)
	}

	return nil
}

func main() {
	departmentCount, err := readInput()
	if err != nil {
		fmt.Printf("error reading department count: %v\n", err)

		return
	}

	for range departmentCount {
		if err = processDepartment(); err != nil {
			fmt.Printf("error processing department: %v\n", err)

			return
		}
	}
}
