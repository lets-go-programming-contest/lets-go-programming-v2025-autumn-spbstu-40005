package main

import (
	"container/heap"
	"fmt"

	"gordey.shapkov/task-2-2/internal/intheap"
)

func main() {
	var amount int
	if _, err := fmt.Scan(&amount); err != nil {
		return
	}

	dishes := &intheap.IntHeap{}

	for range amount {
		var pref int
		if _, err := fmt.Scan(&pref); err != nil {
			return
		}

		heap.Push(dishes, pref)
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		return
	}

	result := findDish(dishes, number)

	fmt.Println(result)
}

func findDish(dishes *intheap.IntHeap, number int) int {
	var value int

	for range dishes.Len() - number + 1 {
		x := heap.Pop(dishes)
		value, _ = x.(int)
	}

	return value
}
