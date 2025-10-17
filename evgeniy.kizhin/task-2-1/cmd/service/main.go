package main

import "errors"
import "fmt"

const minTempLimit = 15
const maxTempLimit = 30

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
		return errors.New("invalid operator")
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
	var operator string
	var temp int
	if _, err := fmt.Scan(&operator, &temp); err != nil {
		return "", 0, err
	}
	return operator, temp, nil
}

func processDepartment(workers_count int) {
	controller := NewTempController()

	for range workers_count {
		operator, temp, err := readWorkerInput()
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		if err := controller.ApplyConstraint(operator, temp); err != nil {
			fmt.Println("Invalid operator:", err)
			return
		}

		fmt.Println(controller.GetOptimalTemp())
	}
}

func main() {
	var departments int
	var workers int
	if _, err := fmt.Scan(&departments); err != nil {
		fmt.Println("Error reading number of departments:", err)
		return
	}
	if _, err := fmt.Scan(&workers); err != nil {
		fmt.Println("Error reading number of workers:", err)
		return
	}

	for range departments {
		processDepartment(workers)
	}
}
