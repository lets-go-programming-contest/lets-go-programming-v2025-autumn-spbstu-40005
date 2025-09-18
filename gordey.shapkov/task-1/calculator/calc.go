package calculator

import "errors"

var (
	errDivisionByZero   = errors.New("Division by zero")
	errInvalidOperation = errors.New("Invalid operation")
)

func Calc(a int, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errDivisionByZero
		} else {
			return a / b, nil
		}
	default:
		return 0, errInvalidOperation
	}
}
