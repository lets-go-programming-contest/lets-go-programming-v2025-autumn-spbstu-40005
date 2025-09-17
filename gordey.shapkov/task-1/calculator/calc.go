package calculator

import "errors"

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
        return 0, errors.New("Division by zero")
      } else {
        return a / b, nil
      }
    default:
      return 0, errors.New("Invalid operation")
  }
}
