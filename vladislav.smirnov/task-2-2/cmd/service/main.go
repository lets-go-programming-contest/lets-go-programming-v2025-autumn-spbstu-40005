package main

import (
	"container/heap"
	"fmt"

	"github.com/smirnov-vladislav/task-2-2/internal/intheap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	var dishesCount int

	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Printf("error scanning count of dishes: %v\n", err)

		return
	}

	dishes := &intheap.IntHeap{}
	heap.Init(dishes)

	for range dishesCount {
		var value int
		if _, err := fmt.Scan(&value); err != nil {
			fmt.Printf("error scanning value of dish: %v\n", err)

			return
		}

		heap.Push(dishes, value)
	}

	var wished int

	if _, err := fmt.Scan(&wished); err != nil {
		fmt.Printf("error scanning wished dish: %v\n", err)

		return
	}

	if wished <= 0 || wished > dishes.Len() {
		fmt.Println("error: wished dish number is out of range")

		return
	}

	for range wished - 1 {
		heap.Pop(dishes)
	}

	bestDish := heap.Pop(dishes)
	fmt.Println(bestDish)
}
