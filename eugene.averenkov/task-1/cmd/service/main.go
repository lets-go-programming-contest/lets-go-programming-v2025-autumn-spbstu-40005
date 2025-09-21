package main

import "fmt"

func main() {
	var a int
	var b int
	var oper string
	_, er1 := fmt.Scan(&a)
	_, er2 := fmt.Scan(&b)
	_, er3 := fmt.Scan(&oper)
	if er1 != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if er2 != nil {
		fmt.Println("Invalid second operand")
		return
	}
	if er3 != nil {
		fmt.Println("Invalid operation")
		return
	}
	switch oper {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(a / b)
	default:
		fmt.Println("Invalid command")
	}
}
