package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dmitei/task-2-1/internal/climate"
)

const expectedPartsCount = 2

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	numDepartments, err := readInt(scanner, "number of departments")
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

func readInt(scanner *bufio.Scanner, what string) (int, error) {
	if !scanner.Scan() {
		return 0, fmt.Errorf("Error: failed to read %s", what)
	}

	value, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return 0, fmt.Errorf("Error: invalid %s", what)
	}

	return value, nil
}

func processDepartment(scanner *bufio.Scanner) error {
	numEmployees, err := readInt(scanner, "number of employees")
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
		return fmt.Errorf("Error: failed to read employee data")
	}

	parts := strings.Fields(strings.TrimSpace(scanner.Text()))
	if len(parts) != expectedPartsCount {
		return fmt.Errorf("Error: invalid input format")
	}

	operator := parts[0]

	temp, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Error: invalid temperature value")
	}

	if err := controller.AddConstraint(operator, temp); err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	fmt.Println(controller.ComfortTemp())

	return nil
}
