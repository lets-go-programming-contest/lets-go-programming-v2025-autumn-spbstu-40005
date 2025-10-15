package minheap

import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func FindKthLargest(ratings []int, preferenceOrder int) int {
	h := &MinHeap{}
	heap.Init(h)

	for _, rating := range ratings {
		if h.Len() < preferenceOrder {
			heap.Push(h, rating)
		} else if rating > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, rating)
		}
	}

	return (*h)[0]
}
