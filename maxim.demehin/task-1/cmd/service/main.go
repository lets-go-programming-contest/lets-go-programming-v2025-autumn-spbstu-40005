package main

import (
	"fmt"
)

func main() {
	var op1, op2 int
	_, err1 := fmt.Scan(&op1)
	_, err2 := fmt.Scan(&op2)

	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operation string
	_, _ = fmt.Scan(&operation)

	switch operation {
	case "+":
		fmt.Println(op1 + op2)

	case "-":
		fmt.Println(op1 - op2)

	case "*":
		fmt.Println(op1 * op2)

	case "/":
		if op2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(op1 / op2)

	default:
		fmt.Println("Invalid operation")
	}
}
