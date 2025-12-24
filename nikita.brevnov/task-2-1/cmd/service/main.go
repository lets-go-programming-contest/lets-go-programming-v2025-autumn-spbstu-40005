package main

import (
	"errors"
	"fmt"
)

var (
	ErrOperator = errors.New("invalid operator")
	ErrOperlow  = errors.New("lowTemp cannot be greater than highTemp")
)

const (
	MinTemp = 15
	MaxTemp = 30
)

type climate struct {
	lowTemp  int
	highTemp int
}

func (oc *climate) AdjustSetting(condition string, value int) error {
	switch condition {
	case ">=":
		if value > oc.lowTemp {
			oc.lowTemp = value
		}
	case "<=":
		if value < oc.highTemp {
			oc.highTemp = value
		}
	default:
		return ErrOperator
	}

	return nil
}

func (oc *climate) FindComfortTemp() int {
	if oc.lowTemp > oc.highTemp {
		return -1
	}

	return oc.lowTemp
}

func handleDep() error {
	var employeeCount int
	if _, err := fmt.Scan(&employeeCount); err != nil {
		return fmt.Errorf("error reading employee count: %w", err)
	}

	climateControl := &climate{
		lowTemp:  MinTemp,
		highTemp: MaxTemp,
	}

	for range employeeCount {
		var (
			condition string
			tempValue int
		)

		if _, err := fmt.Scan(&condition, &tempValue); err != nil {
			return fmt.Errorf("error reading employee input: %w", err)
		}

		if err := climateControl.AdjustSetting(condition, tempValue); err != nil {
			return fmt.Errorf("error adjusting climate settings: %w", err)
		}

		comfortTemp := climateControl.FindComfortTemp()
		fmt.Println(comfortTemp)
	}

	return nil
}

func main() {
	var depCount int

	if _, err := fmt.Scan(&depCount); err != nil {
		fmt.Printf("error reading department count: %v\n", err)

		return
	}

	for range depCount {
		if err := handleDep(); err != nil {
			fmt.Printf("error processing department: %v\n", err)

			return
		}
	}
}
