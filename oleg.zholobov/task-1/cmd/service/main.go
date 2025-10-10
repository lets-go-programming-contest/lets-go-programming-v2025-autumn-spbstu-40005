package main

import "fmt"

func main() {
	var firstOperand int
	fmt.Scan(&firstOperand)

	var secondOperand int
	fmt.Scan(&secondOperand)

	var operation string
	fmt.Scan(&operation)

	fmt.Println(firstOperand)
	fmt.Println(secondOperand)
	fmt.Println(operation)

	var result int
	switch operation {
	case "+":
		result = firstOperand + secondOperand
	case "-":
		result = firstOperand - secondOperand
	case "*":
		result = firstOperand * secondOperand
	case "/":
		result = firstOperand / secondOperand
	default:
		return
	}
	fmt.Println(result)
}
