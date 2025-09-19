package main

import (
	"fmt"

	"gordey.shapkov/task-1/calculator"
)

func main() {
	var (
		a, b      int
		operation string
	)
	_, err := fmt.Scan(&a)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&b)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operator")
		return
	}
	res, err := calculator.Calculation(a, b, operation)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
