package main

import (
	"errors"
	"fmt"
)

var ErrInvalidOperator = errors.New("invalid operator")

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange(minTemp, maxTemp int) *TemperatureRange {
	return &TemperatureRange{
		min: minTemp,
		max: maxTemp,
	}
}

func (tr *TemperatureRange) Update(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tr.min {
			tr.min = temp
		}
	case "<=":
		if temp < tr.max {
			tr.max = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tr *TemperatureRange) GetOptimal() int {
	if tr.min <= tr.max {
		return tr.min
	}

	return -1
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Printf("failed to read departments: %v\n", err)

		return
	}

	for range departments {
		var workers int

		_, err = fmt.Scan(&workers)
		if err != nil {
			fmt.Printf("failed to read workers: %v\n", err)

			return
		}

		temperatureRange := NewTemperatureRange(InitialMinTemp, InitialMaxTemp)

		for range workers {
			var (
				operator string
				temp     int
			)

			_, err = fmt.Scan(&operator, &temp)
			if err != nil {
				fmt.Printf("failed to read temp: %v\n", err)

				return
			}

			err = temperatureRange.Update(operator, temp)
			if err != nil {
				fmt.Printf("failed to update temp: %v\n", err)

				continue
			}

			result := temperatureRange.GetOptimal()
			fmt.Println(result)
		}
	}
}
