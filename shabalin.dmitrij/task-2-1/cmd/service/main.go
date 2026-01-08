package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dmitei/task-2-1/internal/climate"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		fmt.Println("Error: failed to read number of departments")
		return
	}

	numDepartments, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Error: invalid department count")
		return
	}

	for dept := 0; dept < numDepartments; dept++ {
		if !scanner.Scan() {
			fmt.Println("Error: failed to read number of employees")
			return
		}

		numEmployees, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Println("Error: invalid employee count")
			return
		}

		controller := climate.NewController()

		for emp := 0; emp < numEmployees; emp++ {
			if !scanner.Scan() {
				fmt.Println("Error: failed to read employee data")
				return
			}

			parts := strings.Fields(strings.TrimSpace(scanner.Text()))
			if len(parts) != 2 {
				fmt.Println("Error: invalid input format")
				return
			}

			operator := parts[0]
			temp, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid temperature value")
				return
			}

			if err := controller.AddConstraint(operator, temp); err != nil {
				fmt.Println("Error: " + err.Error())
				return
			}

			fmt.Println(controller.ComfortTemp())
		}
	}
}
