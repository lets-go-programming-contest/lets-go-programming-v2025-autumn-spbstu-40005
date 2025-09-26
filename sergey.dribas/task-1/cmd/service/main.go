package main

import (
	"fmt"
)

func main() {
	var (
		leftOperand, rightOperand int
		operation                 string
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
