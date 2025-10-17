package main

import (
	"errors"
	"fmt"
)

const (
	MinTemp = 15
	MaxTemp = 30
)

var ErrInvalidInput = errors.New("invalid input")

type ClimateController struct {
	minTemp int
	maxTemp int
}

func NewClimateController(minTemp, maxTemp int) *ClimateController {
	return &ClimateController{
		minTemp: minTemp,
		maxTemp: maxTemp,
	}
}

func (cc *ClimateController) UpdateBounds(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > cc.minTemp {
			cc.minTemp = temp
		}
	case "<=":
		if temp < cc.maxTemp {
			cc.maxTemp = temp
		}
	default:
		return ErrInvalidInput
	}

	return nil
}

func (cc *ClimateController) GetOptimalTemp() int {
	if cc.minTemp <= cc.maxTemp {
		return cc.minTemp
	}

	return -1
}

func main() {
	var departmentCount int

	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Printf("error reading department count: %v\n", err)

		return
	}

	for range departmentCount {
		var employeeCount int
		if _, err := fmt.Scan(&employeeCount); err != nil {
			fmt.Printf("error reading employee count: %v\n", err)

			return
		}

		controller := NewClimateController(MinTemp, MaxTemp)

		for range employeeCount {
			var (
				operator string
				temp     int
			)

			if _, err := fmt.Scan(&operator, &temp); err != nil {
				fmt.Printf("error reading employee input: %v\n", err)

				return
			}

			if err := controller.UpdateBounds(operator, temp); err != nil {
				fmt.Printf("error updating bounds: %v\n", err)

				return
			}

			result := controller.GetOptimalTemp()
			fmt.Println(result)
		}
	}
}
