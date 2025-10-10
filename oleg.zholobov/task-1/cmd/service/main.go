package main

import "fmt"

func main() {
	var firstOperand int
	_, err := fmt.Scanln(&firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	var secondOperand int
	_, err = fmt.Scanln(&secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var operation string
	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	var result int
	switch operation {
	case "+":
		result = firstOperand + secondOperand
	case "-":
		result = firstOperand - secondOperand
	case "*":
		result = firstOperand * secondOperand
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = firstOperand / secondOperand
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(result)
}
