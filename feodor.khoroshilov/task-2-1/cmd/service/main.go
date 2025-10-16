package main

import (
	"fmt"
	"errors"
)
const (
	InitialMinLimit = 15
	InitialMaxLimit = 30
)

type RangeOfTemperature struct {
	minTemp int
	maxTemp int
}

var ErrorInvalidOperator = errors.New("InvalidOperator")

func NewRangeOfTemperature() *RangeOfTemperature {
	return &RangeOfTemperature{
		minTemp: InitialMinLimit,
		maxTemp: InitialMaxLimit,
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
		return ErrorInvalidOperator
	}
	return nil
}

func processDepartment() error {
	var workers int
	if _, err := fmt.Scan(&workers); err != nil {
		return fmt.Errorf("error when reading number of workers: %w", err)
	}

	tempRange := NewRangeOfTemperature()
	for range workers {
		var operator string
		var temp int

		if _, err := fmt.Scan(&operator, &temp); err != nil {
			return fmt.Errorf("error when reading worker input: %w", err)
		}

		if err := tempRange.SetTemperature(operator, temp); err != nil {
			return fmt.Errorf("error when setting temperature: %w", err)
		}

		result := tempRange.GetOptionalTemperature()
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

	for i := 0; i < numberofdepartments; i++ {
		if err := processDepartment(); err != nil {
			fmt.Printf("error when processing department: %v\n", err)
			return
		}
	}
}
