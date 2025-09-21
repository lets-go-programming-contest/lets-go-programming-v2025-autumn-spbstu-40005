package main

import "fmt"

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
		}
	}
}
