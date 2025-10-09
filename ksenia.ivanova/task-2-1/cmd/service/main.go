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

type Department struct {
	Employees int
	Requests  []Request
}

type Request struct {
	Operator string
	Temp     int
}

func readInt() (int, error) {
	var value int

	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, ErrInvalidInput
	}

	return value, nil
}

func readEmployeeRequest() (Request, error) {
	var operator string
	var temp int

	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return Request{}, ErrInvalidInput
	}

	return Request{Operator: operator, Temp: temp}, nil
}

func readDepartment(employeeCount int) (Department, error) {
	dept := Department{
		Employees: employeeCount,
		Requests:  []Request{},
	}

	for i := range employeeCount {
		req, err := readEmployeeRequest()
		if err != nil {
			return Department{}, ErrInvalidRead
		}

		dept.Requests = append(dept.Requests, req)
	}

	return dept, nil
}

func readInput() ([]Department, error) {
	departmentCount, err := readInt()
	if err != nil {
		return nil, err
	}

	employeeCount, err := readInt()
	if err != nil {
		return nil, err
	}

	departments := make([]Department, 0, departmentCount)

	for i := range departmentCount {
		dept, err := readDepartment(employeeCount)
		if err != nil {
			return nil, ErrInvalidRead
		}

		departments = append(departments, dept)
	}

	return departments, nil
}

func updateTemperatureRange(minTemp, maxTemp int, req Request) (int, int, bool) {
	newMin, newMax := minTemp, maxTemp

	switch req.Operator {
	case ">=":
		if req.Temp > newMin {
			newMin = req.Temp
		}
	case "<=":
		if req.Temp < newMax {
			newMax = req.Temp
		}
	}

	if newMin <= newMax {
		return newMin, newMax, true
	}

	return newMin, newMax, false
}

func processDepartmentRequests(dept Department) []string {
	results := make([]string, 0, len(dept.Requests))
	minTemp, maxTemp := 15, 30
	valid := true

	for _, req := range dept.Requests {
		if !valid {
			results = append(results, "-1")
			continue
		}

		var ok bool
		minTemp, maxTemp, ok = updateTemperatureRange(minTemp, maxTemp, req)

		if ok {
			results = append(results, strconv.Itoa(minTemp))
		} else {
			results = append(results, "-1")
			valid = false
		}
	}

	return results
}

func collectAllResults(departments []Department) []string {
	allResults := make([]string, 0)

	for _, dept := range departments {
		deptResults := processDepartmentRequests(dept)
		allResults = append(allResults, deptResults...)
	}

	return allResults
}

func printResults(results []string) {
	for _, result := range results {
		fmt.Println(result)
	}
}

func main() {
	departments, err := readInput()
	if err != nil {
		fmt.Println("Invalid read")
		return
	}

	results := collectAllResults(departments)
	printResults(results)
}
