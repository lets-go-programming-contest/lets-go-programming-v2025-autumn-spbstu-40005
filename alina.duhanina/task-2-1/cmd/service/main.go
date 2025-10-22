package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidRead  = errors.New("invalid read")
)

type Operator string

const (
	OperatorGreaterOrEqual Operator = ">="
	OperatorLessOrEqual    Operator = "<="
)

type Request struct {
	Operator Operator
	Temp     int
}

func readInt() (int, error) {
	var value int

	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer input: %w", err)
	}

	return value, nil
}

func readEmployeeRequest() (Request, error) {
	var (
		operator string
		temp     int
	)

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return Request{}, fmt.Errorf("failed to read operator and temperature: %w", err)
	}

	return Request{Operator: Operator(operator), Temp: temp}, nil
}

func updateTemperatureRange(minTemp, maxTemp int, req Request) (int, int, error) {
	switch req.Operator {
	case OperatorGreaterOrEqual:
		if req.Temp > minTemp {
			minTemp = req.Temp
		}
	case OperatorLessOrEqual:
		if req.Temp < maxTemp {
			maxTemp = req.Temp
		}
	default:

		return minTemp, maxTemp, ErrInvalidInput
	}

	return minTemp, maxTemp, nil
}

func getTemperatureResult(minTemp, maxTemp int) string {
	if minTemp <= maxTemp {
		return strconv.Itoa(minTemp)
	}

	return "-1"
}

func processDepartmentRequests(employeeCount int) error {
	minTemp, maxTemp := 15, 30

	for range employeeCount {
		req, err := readEmployeeRequest()
		if err != nil {
			return fmt.Errorf("invalid input: employee: %w", err)
		}

		minTemp, maxTemp, err = updateTemperatureRange(minTemp, maxTemp, req)
		if err != nil {
			return fmt.Errorf("process employee request: %w", err)
		}
		result := getTemperatureResult(minTemp, maxTemp)
		fmt.Println(result)
	}

	return nil
}

func readAndProcessInput() error {
	departmentCount, err := readInt()
	if err != nil {
		return err
	}

	for range departmentCount {
		employeeCount, err := readInt()
		if err != nil {
			return err
		}

		err = processDepartmentRequests(employeeCount)
		if err != nil {
			return fmt.Errorf("invalid input: department: %w", err)
		}
	}

	return nil
}

func main() {
	err := readAndProcessInput()
	if err != nil {
		fmt.Printf("input processing error: %w\n", err)

		return
	}

}
