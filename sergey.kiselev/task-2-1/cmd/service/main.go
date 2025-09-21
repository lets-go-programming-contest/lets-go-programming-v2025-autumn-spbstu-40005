package main

import "fmt"

const (
	maxTemperature = 30
	minTemperature = 15
)

func main() {
	var countDepartaments int

	if _, err := fmt.Scan(&countDepartaments); err != nil || countDepartaments < 1 {
		fmt.Println("Incorrect countDepartaments")

		return
	}

	for range countDepartaments {
		var countEmployees int

		if _, err := fmt.Scan(&countEmployees); err != nil || countEmployees < 1 {
			fmt.Println("Incorrect countEmployees")

			return
		}

		maxT := maxTemperature
		minT := minTemperature

		for range countEmployees {
			var (
				operator    string
				temperature int
			)

			if _, err := fmt.Scan(&operator, &temperature); err != nil {
				fmt.Println("Invalid input for operator and temperature")

				return
			}

			switch operator {
			case "<=":
				if temperature < maxT {
					maxT = temperature
				}
			case ">=":
				if temperature > minT {
					minT = temperature
				}
			default:
				fmt.Println("Incorrect operator")

				return
			}

			if minT <= maxT {
				fmt.Println(minT)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
