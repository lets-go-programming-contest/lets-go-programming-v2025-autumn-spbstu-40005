package main

import (
	"fmt"
)

func main() {
	var (
		leftOperand, rightOperand int
		operation                 string
	)
	n, err := fmt.Scan(&leftOperand, &rightOperand, &operation)
	if err != nil {
		switch n {
		case 0:
			fmt.Println("Invalid first operand")
		case 1:
			fmt.Println("Invalid second operand")
		case 2:
			fmt.Println("Invalid operation")
		}
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
