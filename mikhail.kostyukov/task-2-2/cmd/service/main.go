package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/KostyukovMichael/lets-go-programming-v2025-autumn-spbstu-40005/task-2-2/internal/intheap"
)

var (
	errConvertToInt = errors.New("cannot convert to int")
	errHeapIsEmpty  = errors.New("cannot pop from an empty heap")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %v\n", err)
		}
	}()

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

	for range wished - 1 {
		heap.Pop(dishesHeap)
	}

	popped := heap.Pop(dishesHeap)
	if popped == nil {
		fmt.Printf("while popping last element from the heap: %v\n", errHeapIsEmpty)

		return
	}

	result, ok := popped.(int)
	if !ok {
		fmt.Printf("while converting last element to int: %v\n", errConvertToInt)
	}

	fmt.Println(result)
}
