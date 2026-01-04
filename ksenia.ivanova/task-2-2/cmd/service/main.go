package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

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

func main() {
	var n, k int

	fmt.Scan(&n)

	ratings := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Scan(&ratings[i])
	}

	fmt.Scan(&k)

	if k < 1 || k > n {
		return
	}

	h := &MinHeap{}
	heap.Init(h)

	for i := 0; i < k; i++ {
		heap.Push(h, ratings[i])
	}

	for i := k; i < n; i++ {
		if ratings[i] > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, ratings[i])
		}
	}

	result := (*h)[0]
	fmt.Println(result)
}
