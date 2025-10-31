package main

import (
	"errors"
	"fmt"
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
		return ErrInvalidInput
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
		fmt.Printf("input processing error: %v\n", err)

		return
	}
}
