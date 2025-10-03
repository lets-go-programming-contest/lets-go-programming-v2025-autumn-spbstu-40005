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

	for range departmentSize {
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
	var size int
	_, err := fmt.Scan(&size)

	if err != nil {
		return
	}
	
	for range size {
		if !processDepartment() {
			return
		}
	}
}
