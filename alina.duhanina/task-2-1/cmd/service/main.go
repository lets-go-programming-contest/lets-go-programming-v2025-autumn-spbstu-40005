package main

import (
	"errors"
	"fmt"
)

var ErrInvalidOperator = errors.New("invalid operator")

type Operator string

const (
	OperatorGreaterOrEqual Operator = ">="
	OperatorLessOrEqual    Operator = "<="
)

const (
	InitialMinTemp = 15
	InitialMaxTemp = 30
)

type Request struct {
	Operator Operator
	Temp     int
}

type TemperatureRange struct {
	minTemp int
	maxTemp int
}

func NewTemperatureRange(minTemp, maxTemp int) *TemperatureRange {
	return &TemperatureRange{
		minTemp: minTemp,
		maxTemp: maxTemp,
	}
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

func (tr *TemperatureRange) updateTemperatureRange(req Request) error {
	switch req.Operator {
	case OperatorGreaterOrEqual:
		if req.Temp > tr.minTemp {
			tr.minTemp = req.Temp
		}
	case OperatorLessOrEqual:
		if req.Temp < tr.maxTemp {
			tr.maxTemp = req.Temp
		}
	default:
		return ErrInvalidOperator
	}

	return nil
}

func (tr *TemperatureRange) getTemperatureResult() int {
	if tr.minTemp <= tr.maxTemp {
		return tr.minTemp
	}

	return -1
}

func processDepartmentRequests(employeeCount int) error {
	tempRange := NewTemperatureRange(InitialMinTemp, InitialMaxTemp)

	for range employeeCount {
		req, err := readEmployeeRequest()
		if err != nil {
			return fmt.Errorf("invalid input: employee: %w", err)
		}

		err = tempRange.updateTemperatureRange(req)
		if err != nil {
			return fmt.Errorf("process employee request: %w", err)
		}

		result := tempRange.getTemperatureResult()
		fmt.Println(result)
	}

	return nil
}

func main() {
	departmentCount, err := readInt()
	if err != nil {
		fmt.Printf("failed to read department count: %v\n", err)

		return
	}

	for range departmentCount {
		employeeCount, err := readInt()
		if err != nil {
			fmt.Printf("failed to read employee count for department: %v\n", err)

			return
		}

		err = processDepartmentRequests(employeeCount)
		if err != nil {
			fmt.Printf("failed to process department requests: %v\n", err)

			return
		}
	}
}
