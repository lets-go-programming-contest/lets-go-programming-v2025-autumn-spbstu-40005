package main

import (
	"container/heap"
	"fmt"
	"gordey.shapkov/task-2-2/internal/IntHeap"
)

func main() {
	var amount int
	if _, err := fmt.Scan(&amount); err != nil {
		return
	}

	dishes := &IntHeap.IntHeap{}
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

	var result int
	for range amount - number + 1 {
		popped := heap.Pop(dishes)
		result = popped.(int)
	}
	fmt.Println(result)
}
