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

type Request struct {
	Operator string
	Temp     int
}

func readInt() (int, error) {
	var value int

	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrInvalidInput, err)
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
		return Request{}, fmt.Errorf("%s: %w", ErrInvalidInput, err)
	}

	return Request{Operator: operator, Temp: temp}, nil
}

func updateTemperatureRange(minTemp, maxTemp int, req Request) (int, int) {
	switch req.Operator {
	case ">=":
		if req.Temp > minTemp {
			minTemp = req.Temp
		}
	case "<=":
		if req.Temp < maxTemp {
			maxTemp = req.Temp
		}
	}

	return minTemp, maxTemp
}

func getTemperatureResult(minTemp, maxTemp int) string {
	if minTemp <= maxTemp {
		return strconv.Itoa(minTemp)
	}

	return "-1"
}

func processDepartmentRequests(employeeCount int) ([]string, error) {
	results := make([]string, 0, employeeCount)
	minTemp, maxTemp := 15, 30

	for range employeeCount {
		req, err := readEmployeeRequest()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrInvalidRead, err)
		}

		minTemp, maxTemp = updateTemperatureRange(minTemp, maxTemp, req)

		results = append(results, getTemperatureResult(minTemp, maxTemp))
	}

	return results, nil
}

func readAndProcessInput() ([]string, error) {
	departmentCount, err := readInt()
	if err != nil {
		return nil, err
	}

	allResults := make([]string, 0)

	for range departmentCount {
		employeeCount, err := readInt()
		if err != nil {
			return nil, err
		}

		deptResults, err := processDepartmentRequests(employeeCount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ErrInvalidRead, err)
		}

		allResults = append(allResults, deptResults...)
	}

	return allResults, nil
}

func printResults(results []string) {
	for _, result := range results {
		fmt.Println(result)
	}
}

func main() {
	allResults, err := readAndProcessInput()
	if err != nil {
		fmt.Printf("Invalid read: %v", err)

		return
	}

	printResults(allResults)
}
