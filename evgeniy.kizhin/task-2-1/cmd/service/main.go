package main

import (
	"errors"
	"fmt"
)

const (
	minTempLimit = 15
	maxTempLimit = 30
)

var (
	ErrInvalidOperator = errors.New("invalid operator")
	ErrReadingWorkers  = errors.New("error reading number of workers")
	ErrReadingInput    = errors.New("error reading worker input")
	ErrUpdatingTemp    = errors.New("error updating temperature manager")
)

type TempController struct {
	minTemp int
	maxTemp int
}

func NewTempController() *TempController {
	return &TempController{
		minTemp: minTempLimit,
		maxTemp: maxTempLimit,
	}
}

func (tc *TempController) ApplyConstraint(operator string, temp int) error {
	switch operator {
	case ">=":
		if temp > tc.minTemp {
			tc.minTemp = temp
		}
	case "<=":
		if temp < tc.maxTemp {
			tc.maxTemp = temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tc *TempController) GetOptimalTemp() int {
	if tc.minTemp > tc.maxTemp {
		return -1
	}

	return tc.minTemp
}

func readWorkerInput() (string, int, error) {
	var (
		operator string
		temp     int
	)

	if _, err := fmt.Scan(&operator, &temp); err != nil {
		return "", 0, ErrReadingInput
	}

	return operator, temp, nil
}

func processDepartment() error {
	var workers int
	if _, err := fmt.Scan(&workers); err != nil {
		return ErrReadingWorkers
	}

	controller := NewTempController()

	for range workers {
		operator, temp, err := readWorkerInput()
		if err != nil {
			return err
		}

		if err := controller.ApplyConstraint(operator, temp); err != nil {
			return ErrUpdatingTemp
		}

		fmt.Println(controller.GetOptimalTemp())
	}

	return nil
}

func main() {
	var departments int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println(ErrReadingWorkers)

		return
	}

	for range departments {
		if err := processDepartment(); err != nil {
			fmt.Println(err)

			return
		}
	}
}
