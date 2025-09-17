package main

import (
	"fmt"
	"os"
)

func main() {
	var op1, op2 int
	_, err1 := fmt.Scan(&op1)
	_, err2 := fmt.Scan(&op2)

	if err1 != nil {
		fmt.Println("Invalid first operand")
		os.Exit(1)
	}
	if err2 != nil {
		fmt.Println("Invalid second operand")
		os.Exit(1)
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
			os.Exit(1)
		}
		fmt.Println(op1 / op2)

	default:
		fmt.Println("Invalid operation")
	}

}
