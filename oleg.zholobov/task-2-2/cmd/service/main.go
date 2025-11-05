package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (heap *MinHeap) Len() int {
	return len(*heap)
}

func (heap *MinHeap) Less(i, j int) bool {
	return (*heap)[i] < (*heap)[j]
}

func (heap *MinHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *MinHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		panic("MinHeap: expected int value")
	}

	*heap = append(*heap, v)
}

func (heap *MinHeap) Pop() any {
	old := *heap

	length := len(old)
	if length == 0 {
		return nil
	}

	x := old[length-1]
	*heap = old[0 : length-1]

	return x
}

func (heap *MinHeap) Peek() (int, bool) {
	if len(*heap) == 0 {
		return 0, false
	}

	return (*heap)[0], true
}

func main() {
	var amountOfDishes int
	if _, err := fmt.Scan(&amountOfDishes); err != nil {
		fmt.Println("Error: invalid format for number of dishes")

		return
	}

	if amountOfDishes <= 0 {
		fmt.Println("Error: number of dishes must be a positive number")

		return
	}

	dishes := make([]int, amountOfDishes)
	for i := range dishes {
		if _, err := fmt.Scan(&dishes[i]); err != nil {
			fmt.Println("Error: invalid format for dish rating")

			return
		}
	}

	var dishNumber int
	if _, err := fmt.Scan(&dishNumber); err != nil {
		fmt.Println("Error: invalid format for dish number")

		return
	}

	if dishNumber <= 0 {
		fmt.Println("Error: dish number must be a positive number")

		return
	}

	if dishNumber > amountOfDishes {
		fmt.Printf("Error: dish number (%d) cannot exceed number of dishes (%d)\n", dishNumber, amountOfDishes)

		return
	}

	dishesHeap := &MinHeap{}
	heap.Init(dishesHeap)

	for _, val := range dishes {
		heap.Push(dishesHeap, val)

		if dishesHeap.Len() > dishNumber {
			heap.Pop(dishesHeap)
		}
	}

	if result, ok := dishesHeap.Peek(); ok {
		fmt.Println(result)
	} else {
		fmt.Println("No dish available")
	}
}
