package main

import "fmt"

func main() {

	var (
		leftOperand  int
		rightOperand int
		operation    string
		result       int
	)

	n, err := fmt.Scan(&leftOperand, &rightOperand)

	if err != nil {
		if n == 0 {
			fmt.Println("Invalid first operand")
		} else if n == 1 {
			fmt.Println("Invalid second operand")
		}
		return
	}

	fmt.Scan(&operation)

	switch operation {
	case "+":
		result = leftOperand + rightOperand
	case "-":
		result = leftOperand - rightOperand
	case "*":
		result = leftOperand * rightOperand
	case "/":
		if rightOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = leftOperand / rightOperand
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(result)
}
