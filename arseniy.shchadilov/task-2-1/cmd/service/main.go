package main

import (
	"errors"
	"fmt"
)

const (
	GlobalMinTemp = 15
	GlobalMaxTemp = 30
)

var ErrInvalidOperator = errors.New("update failed: invalid operator")

type TemperatureManager struct {
	currentMinTemp int
	currentMaxTemp int
}

func NewTemperatureManager(minTemp, maxTemp int) *TemperatureManager {
	return &TemperatureManager{
		currentMinTemp: minTemp,
		currentMaxTemp: maxTemp,
	}
}

func (tm *TemperatureManager) Update(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tm.currentMinTemp {
			tm.currentMinTemp = temp
		}
	case "<=":
		if temp < tm.currentMaxTemp {
			tm.currentMaxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tm *TemperatureManager) GetOptimalTemp() int {
	if tm.currentMinTemp > tm.currentMaxTemp {
		return -1
	}

	return tm.currentMinTemp
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("error: failed to read number of departments")

		return
	}

	for range departments {
		var workers int

		_, err := fmt.Scan(&workers)
		if err != nil {
			fmt.Println("error: failed to read number of workers for department")

			return
		}

		manager := NewTemperatureManager(GlobalMinTemp, GlobalMaxTemp)

		for range workers {
			var (
				operator string
				temp     int
			)

			_, err := fmt.Scan(&operator, &temp)
			if err != nil {
				fmt.Println("error: failed to read worker")

				return
			}

			err = manager.Update(operator, temp)
			if err != nil {
				fmt.Println("error: failed to update temperature manager for worker")

				return
			}

			fmt.Println(manager.GetOptimalTemp())
		}
	}
}
