package minheap

import "container/heap"

type MinHeap []int

func (heap *MinHeap) Len() int {
	return len(*heap)
}

func (heap *MinHeap) Less(i, j int) bool {
	if i >= len(*heap) || j >= len(*heap) {
		panic("array index out of bounds")
	}

	return (*heap)[i] < (*heap)[j]
}

func (heap *MinHeap) Swap(i, j int) {
	if i >= len(*heap) || j >= len(*heap) {
		panic("array index out of bounds")
	}

	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("invalid argument")
	}

	*heap = append(*heap, value)
}

func (heap *MinHeap) Pop() interface{} {
	if len(*heap) == 0 {
		return nil
	}

	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

func FindKthLargest(ratings []int, preferenceOrder int) int {
	heapInstance := &MinHeap{}
	heap.Init(heapInstance)

	for _, rating := range ratings {
		if heapInstance.Len() < preferenceOrder {
			heap.Push(heapInstance, rating)
		} else if rating > (*heapInstance)[0] {
			heap.Pop(heapInstance)
			heap.Push(heapInstance, rating)
		}
	}

	return (*heapInstance)[0]
}
