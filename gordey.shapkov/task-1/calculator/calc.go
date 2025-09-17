package calculator

import "fmt"

func Calc(a float64, b float64, op string) {
  switch op {
    case "+":
      fmt.Println(a + b)
    case "-":
      fmt.Println(a - b)
    case "*":
      fmt.Println(a * b)
    case "/":
      if b == 0 {
        fmt.Println("Division by zero!")
      } else {
        fmt.Println(a / b)
      }
    default:
      fmt.Println("Invalid operation")
  }
}
