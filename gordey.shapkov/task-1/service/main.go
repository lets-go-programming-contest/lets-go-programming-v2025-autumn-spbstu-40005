package main

import "fmt"
import "first/calculator"

func main() {
  var (
      a, b float64
      op string
      )
  _, err := fmt.Scan(&a)
  if err != nil {
    fmt.Println("Invalid first operand")
    return
  }
  _, err = fmt.Scan(&b)
  if err != nil {
    fmt.Println("Invalid second operand")
    return
  }
  fmt.Scan(&op)
  calculator.Calc(a, b, op)
}
