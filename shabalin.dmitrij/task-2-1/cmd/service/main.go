package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dmitei/task-2-1/internal/climate"
)

const expectedPartsCount = 2

var (
	errReadDepartments    = errors.New("failed to read number of departments")
	errInvalidDepartments = errors.New("invalid department count")

	errReadEmployees    = errors.New("failed to read number of employees")
	errInvalidEmployees = errors.New("invalid employee count")

	errReadEmployeeData   = errors.New("failed to read employee data")
	errInvalidInputFormat = errors.New("invalid input format")
	errInvalidTemperature = errors.New("invalid temperature value")
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	numDepartments, err := readDepartmentsCount(scanner)
	if err != nil {
		fmt.Println(err)

		return
	}

	for range numDepartments {
		if err := processDepartment(scanner); err != nil {
			fmt.Println(err)

			return
		}
	}
}

func readDepartmentsCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errReadDepartments
	}

	value, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return 0, errInvalidDepartments
	}

	return value, nil
}

func readEmployeesCount(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errReadEmployees
	}

	value, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return 0, errInvalidEmployees
	}

	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	numEmployees, err := readEmployeesCount(scanner)
	if err != nil {
		return err
	}

	controller := climate.NewController()

	for range numEmployees {
		if err := processEmployee(scanner, controller); err != nil {
			return err
		}
	}

	return nil
}

func processEmployee(scanner *bufio.Scanner, controller *climate.Controller) error {
	if !scanner.Scan() {
		return errReadEmployeeData
	}

	parts := strings.Fields(strings.TrimSpace(scanner.Text()))
	if len(parts) != expectedPartsCount {
		return errInvalidInputFormat
	}

	operator := parts[0]

	temp, err := strconv.Atoi(parts[1])
	if err != nil {
		return errInvalidTemperature
	}

	if err := controller.AddConstraint(operator, temp); err != nil {
		return fmt.Errorf("failed to apply constraint: %w", err)
	}

	fmt.Println(controller.ComfortTemp())

	return nil
}
