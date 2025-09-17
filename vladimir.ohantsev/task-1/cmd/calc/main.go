package main

import (
	"errors"
	"fmt"
	"os"
)

func plus(x, y int) (int, error) {
	return x + y, nil
}

func minus(x, y int) (int, error) {
	return x - y, nil
}

func mult(x, y int) (int, error) {
	return x * y, nil
}

func div(x, y int) (int, error) {
	if y == 0 {
		return 0, errors.New("Division by zero")
	}
	return x / y, nil
}

func main() {
	m := map[string]func(x, y int) (int, error){
		"+": plus,
		"-": minus,
		"*": mult,
		"/": div,
	}
	var x, y int
	var op string
	if _, err := fmt.Scan(&x); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid first operand")
		return
	}
	if _, err := fmt.Scan(&y); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid second operand")
		return
	}
	_, _ = fmt.Scan(&op)
	if f, ok := m[op]; ok {
		res, err := f(x, y)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		fmt.Println(res)
	} else {
		fmt.Fprintln(os.Stderr, "Invalid operation")
		return
	}
}
