package minintheap

import "container/heap"

type MinIntHeap []int

func (h MinIntHeap) Len() int {
	return len(h)
}

func (h MinIntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinIntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinIntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}
	*h = append(*h, value)
}

func (h *MinIntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func KthLargest(values []int, k int) int {
	h := &MinIntHeap{}
	heap.Init(h)

	for _, value := range values {
		if h.Len() < k {
			heap.Push(h, value)
		} else if value > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, value)
		}
	}

	return (*h)[0]
}
