package main

import (
	"fmt"

	"alexandra.karnauhova/task-2-1/internal/temperature"
)

const (
	defaultMax = 100
	defaultMin = 0
)

func inputDataDepartment() {
	department := temperature.TempManager{
		Temps:            []temperature.RangeTemp{},
		IdealTemperature: defaultMin,
		Max:              defaultMax,
		Min:              defaultMin,
	}

	var numberEmployee int

	_, err := fmt.Scan(&numberEmployee)
	if err != nil {
		fmt.Println("Invalid number of employee")

		return
	}

	for range numberEmployee {
		var (
			signs  string
			number int
		)

		_, err := fmt.Scan(&signs, &number)
		if err != nil {
			fmt.Println("Invalid temperature")

			return
		}

		err = department.AddTemp(signs, number)
		if err != nil {
			fmt.Println(-1)
		} else {
			fmt.Println(department.IdealTemperature)
		}
	}
}

func main() {
	var numberDepartment int

	_, err := fmt.Scan(&numberDepartment)
	if err != nil {
		fmt.Println("Invalid number of departments")

		return
	}

	for range numberDepartment {
		inputDataDepartment()
	}
}
