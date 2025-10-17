package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	if n <= 0 {
		return
	}

	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&arr[i]); err != nil {
			return
		}
	}

	var k int
	if _, err := fmt.Scan(&k); err != nil {
		return
	}

	if k <= 0 || k > n {
		return
	}

	fmt.Println(n)
	for _, val := range arr {
		fmt.Println(val)
	}
	fmt.Println(k)

}
