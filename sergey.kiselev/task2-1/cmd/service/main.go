package main

import "fmt"

const (
	maxTemperature = 30
	minTemperature = 15
)

func main() {
	var countDepartaments int

	_, err := fmt.Scan(&countDepartaments)
	if err != nil {
		fmt.Println("Incorrect countDepartaments")
		return
	}

	for i := 0; i < countDepartaments; i++ {
		var countEmployees int

		_, err = fmt.Scan(&countEmployees)
		if err != nil {
			fmt.Println("Incorrect countEmployees")
			return
		}

		maxT := maxTemperature
		minT := minTemperature

		for j := 0; j < countEmployees; j++ {
			var (
				operator    string
				temperature int
			)

			_, err = fmt.Scan(&operator, &temperature)
			if err != nil {
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
				fmt.Print("Incorrect operator")
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
