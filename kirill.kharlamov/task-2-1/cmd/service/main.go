package main

import (
	"errors"
	"fmt"
)

const (
	MinAllowed = 15
	MaxAllowed = 30
)

var (
	ErrInvalidOperator = errors.New("invalid operator")
	ErrTempOutOfRange  = errors.New("temperature out of allowed range")
	ErrNoComfortTemp   = errors.New("no comfortable temperature found")
)

type OfficeClimate struct {
	lowBound  int
	highBound int
}

func CreateClimateController() *OfficeClimate {
	return &OfficeClimate{
		lowBound:  MinAllowed,
		highBound: MaxAllowed,
	}
}

func (oc *OfficeClimate) AdjustSetting(condition string, value int) error {
	if value < MinAllowed || value > MaxAllowed {
		return ErrTempOutOfRange
	}

	switch condition {
	case ">=":
		if value > oc.lowBound {
			oc.lowBound = value
		}
	case "<=":
		if value < oc.highBound {
			oc.highBound = value
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (oc *OfficeClimate) FindComfortTemp() (int, error) {
	if oc.lowBound > oc.highBound {
		return 0, ErrNoComfortTemp
	}

	return oc.lowBound, nil
}

func handleDepartment() error {
	var employeeCount int
	if _, err := fmt.Scan(&employeeCount); err != nil {
		return fmt.Errorf("error reading employee count: %w", err)
	}

	climateControl := CreateClimateController()

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

		comfortTemp, err := climateControl.FindComfortTemp()

		if err != nil {
			if errors.Is(err, ErrNoComfortTemp) {
				fmt.Println(-1)
				continue
			}

			return err
		}

		fmt.Println(comfortTemp)
	}

	return nil
}

func main() {
	var departmentCount int
	if _, err := fmt.Scan(&departmentCount); err != nil {
		fmt.Printf("error reading department count: %v\n", err)

		return
	}

	for range departmentCount {
		if err := handleDepartment(); err != nil {
			fmt.Printf("error processing department: %v\n", err)

			continue
		}
	}
}
