package main

import (
	"fmt"
)

func main() {
	var (
		leftOperand, rightOperand int
		operation                 string
	)
	_, err := fmt.Scan(&leftOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&rightOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(leftOperand + rightOperand)
	case "-":
		fmt.Println(leftOperand - rightOperand)
	case "*":
		fmt.Println(leftOperand * rightOperand)
	case "/":
		if rightOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(leftOperand / rightOperand)
	default:
		fmt.Println("Invalid operation")
	}
}
