package main

import (
	"fmt"
	"strconv"
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
		return 0, fmt.Errorf("Invalid input")
	}
	return value, nil
}

func readEmployeeRequest() (Request, error) {
	var operator string
	var temp int
	_, err := fmt.Scan(&operator, &temp)
	if err != nil {
		return Request{}, fmt.Errorf("Invalid input")
	}
	return Request{Operator: operator, Temp: temp}, nil
}

func readDepartment() (Department, error) {
	K, err := readInt()
	if err != nil {
		return Department{}, err
	}
	dept := Department{Employees: K}
	for j := 0; j < K; j++ {
		req, err := readEmployeeRequest()
		if err != nil {
			return Department{}, fmt.Errorf("Invalid read")
		}
		dept.Requests = append(dept.Requests, req)
	}
	return dept, nil
}

func readInput() ([]Department, error) {
	N, err := readInt()
	if err != nil {
		return nil, err
	}
	var departments []Department
	for i := 0; i < N; i++ {
		dept, err := readDepartment()
		if err != nil {
			return nil, fmt.Errorf("Invalid read")
		}
		departments = append(departments, dept)
	}
	return departments, nil
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

func processDepartmentRequests(dept Department) []string {
	var results []string
	minTemp, maxTemp := 15, 30
	for _, req := range dept.Requests {
		minTemp, maxTemp = updateTemperatureRange(minTemp, maxTemp, req)
		results = append(results, getTemperatureResult(minTemp, maxTemp))
	}
	return results
}

func collectAllResults(departments []Department) []string {
	var allResults []string
	for _, dept := range departments {
		allResults = append(allResults, processDepartmentRequests(dept)...)
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
		fmt.Printf("Invalid read")
		return
	}
	results := collectAllResults(departments)
	printResults(results)
}
