package main

import "fmt"

func main() {
	var number1 int
	_, err := fmt.Scan(&number1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	var number2 int
	_, err = fmt.Scan(&number2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}
	var oper string
	_, err = fmt.Scan(&oper)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}
	switch oper {
	case "+":
		fmt.Println(number1 + number2)
	case "-":
		fmt.Println(number1 - number2)
	case "*":
		fmt.Println(number1 * number2)
	case "/":
		if number2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(number1 / number2)
	default:
		fmt.Println("Invalid operation")
	}
}
