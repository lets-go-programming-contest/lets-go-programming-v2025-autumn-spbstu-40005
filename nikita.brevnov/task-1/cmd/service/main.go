package main

import "fmt"

func main() {
	var oper1 int
	var oper2 int
	var operation string
	_, err := fmt.Scan(&oper1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&oper2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println("Invalid operation")
	}
	switch operation {
	case "+":
		fmt.Println(oper1 + oper2)
	case "-":
		fmt.Println(oper1 - oper2)
	case "*":
		fmt.Println(oper1 * oper2)
	case "/":
		if oper2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(oper1 / oper2)
	default:
		fmt.Println("Invalid operation")
	}
}
