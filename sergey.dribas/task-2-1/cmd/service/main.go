package main

import "fmt"

func processDepartment() bool {
	var departmentSize int
	_, err := fmt.Scan(&departmentSize)
	if err != nil {
		return false
	}
	minTemp := 15
	maxTemp := 30
	for i := 0; i < departmentSize; i++ {
		var operand string
		var temp int
		_, err := fmt.Scan(&operand, &temp)
		if err != nil {
			return false
		}
		switch operand {
		case ">=":
			if temp > minTemp {
				minTemp = temp
			}
		case "<=":
			if temp < maxTemp {
				maxTemp = temp
			}
		default:

			return false
		}
		if minTemp > maxTemp {
			fmt.Println("-1")
		} else {
			fmt.Println(minTemp)
		}
	}
	return true
}

func main() {
	var n int
	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}
	for i := 0; i < n; i++ {
		if !processDepartment() {
			return
		}
	}
}
