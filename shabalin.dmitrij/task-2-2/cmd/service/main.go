package main

import (
	"fmt"

	"github.com/dmitei/task-2-2/internal/heap"
)

func main() {
	var n int
	fmt.Scan(&n)

	ratings := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&ratings[i])
	}

	var k int
	fmt.Scan(&k)

	result := heap.FindKthPreferred(ratings, k)
	fmt.Println(result)
}
