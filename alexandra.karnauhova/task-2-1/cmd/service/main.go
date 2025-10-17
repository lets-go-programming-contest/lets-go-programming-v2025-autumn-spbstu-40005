package main

import (
	"fmt"

	"alexandra.karnauhova/task-2-1/internal/temperature"
)

func inputDataDepartment() {
	department := temperature.TempManager{
		Temps:            []temperature.RangeTemp{},
		IdealTemperature: 0,
		Max:              100,
		Min:              0,
	}

	var numberEmployee int

	_, err := fmt.Scan(&numberEmployee)
	if err != nil {
		fmt.Println("Invalid number of employee")
		return
	}

	for i := 0; i < numberEmployee; i++ {
		var signs string
		var number int
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

	for i := 0; i < numberDepartment; i++ {
		inputDataDepartment()
	}
}
