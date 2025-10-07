package main

import (
	"container/heap"
	"fmt"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-2-2/internal/intheap"
)

func main() {
	var dishesCnt int
	if _, err := fmt.Scan(&dishesCnt); err != nil {
		fmt.Printf("invalid input of number of dishes: %v\n", err)

		return
	}

	dishesHeap := &intheap.IntHeap{}
	heap.Init(dishesHeap)

	for range dishesCnt {
		var preference int
		if _, err := fmt.Scan(&preference); err != nil {
			fmt.Printf("invalid input of preference: %v\n", err)

			return
		}

		heap.Push(dishesHeap, preference)
	}

	var wished int
	if _, err := fmt.Scan(&wished); err != nil {
		fmt.Printf("invalid input of wished dish: %v\n", err)

		return
	}

	if wished < 1 || wished > dishesHeap.Len() {
		fmt.Println("that dish does not exist")

		return
	}

	var result int
	for range wished {
		result = heap.Pop(dishesHeap).(int)
	}

	fmt.Println(result)
}
