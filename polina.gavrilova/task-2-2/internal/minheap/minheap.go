package minheap

import "errors"

type MinHeap []int

var ErrHeapEmpty = errors.New("heap is empty")

func (heap *MinHeap) Len() int {
	return len(*heap)
}

func (heap *MinHeap) Less(i, j int) bool {
	if i < 0 || j < 0 || i >= heap.Len() || j >= heap.Len() {
		panic("less index is out of range")
	}

	return (*heap)[i] < (*heap)[j]
}

func (heap *MinHeap) Swap(i, j int) {
	if i < 0 || j < 0 || i >= heap.Len() || j >= heap.Len() {
		panic("swap index is out of range")
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

func (heap *MinHeap) Top() (int, error) {
	if len(*heap) == 0 {
		return 0, ErrHeapEmpty
	}

	return (*heap)[0], nil
}
