package main

import (
	"errors"
	"fmt"
)

const (
	InitialMinLimit = 15
	InitialMaxLimit = 30
)

type RangeOfTemperature struct {
	minTemp int
	maxTemp int
}

var ErrInvalidOperator = errors.New("unsupported operator type")

func NewRangeOfTemperature(minTemp, maxTemp int) *RangeOfTemperature {
	return &RangeOfTemperature{
		minTemp: minTemp,
		maxTemp: maxTemp,
	}
}

func (tmp *RangeOfTemperature) GetOptionalTemperature() int {
	if tmp.minTemp > tmp.maxTemp {
		return -1
	}

	return tmp.minTemp
}

func (tmp *RangeOfTemperature) SetTemperature(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tmp.minTemp {
			tmp.minTemp = temp
		}
	case "<=":
		if temp < tmp.maxTemp {
			tmp.maxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func processDepartment() error {
	var workers int
	if _, err := fmt.Scan(&workers); err != nil {
		return fmt.Errorf("error when reading number of workers: %w", err)
	}

	temprange := NewRangeOfTemperature(InitialMinLimit, InitialMaxLimit)

	for range workers {
		var (
			operator string
			temp     int
		)

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return fmt.Errorf("error when reading worker input: %w", err)
		}

		if err := temprange.SetTemperature(operator, temp); err != nil {
			return fmt.Errorf("error when setting temperature: %w", err)
		}

		result := temprange.GetOptionalTemperature()
		fmt.Println(result)
	}

	return nil
}

func main() {
	var numberofdepartments int
	if _, err := fmt.Scan(&numberofdepartments); err != nil {
		fmt.Printf("error when reading number of departments: %v\n", err)

		return
	}

	for range numberofdepartments {
		if err := processDepartment(); err != nil {
			fmt.Printf("error when processing department: %v\n", err)

			return
		}
	}
}
