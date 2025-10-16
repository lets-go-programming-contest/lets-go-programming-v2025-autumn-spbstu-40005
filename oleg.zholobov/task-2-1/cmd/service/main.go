package main

import "fmt"

func readInput() (int, error) {
	var value int
	_, err := fmt.Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func main() {
	departmentCount, err := readInput()
	if err != nil {
		fmt.Println(-1)
		return
	}

	fmt.Println(departmentCount)
}
