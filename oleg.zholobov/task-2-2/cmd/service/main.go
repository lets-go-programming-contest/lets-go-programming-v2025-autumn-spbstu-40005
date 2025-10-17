package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (heap MinHeap) Len() int           { return len(heap) }
func (heap MinHeap) Less(i, j int) bool { return heap[i] < heap[j] }
func (heap MinHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

func (heap *MinHeap) Push(x any) {
	*heap = append(*heap, x.(int))
}

func (heap *MinHeap) Pop() any {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]
	return x
}

func main() {
	var amountOfDishes int
	if _, err := fmt.Scan(&amountOfDishes); err != nil {
		return
	}

	if amountOfDishes <= 0 {
		return
	}

	dishes := make([]int, amountOfDishes)
	for i := 0; i < amountOfDishes; i++ {
		if _, err := fmt.Scan(&dishes[i]); err != nil {
			return
		}
	}

	var dishNumber int
	if _, err := fmt.Scan(&dishNumber); err != nil {
		return
	}

	if dishNumber <= 0 || dishNumber > amountOfDishes {
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

	fmt.Println((*dishesHeap)[0])
}
