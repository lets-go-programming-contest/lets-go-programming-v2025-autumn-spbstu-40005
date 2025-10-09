package intheap

import (
	"errors"
	"sort"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	val, err := x.(int)
	if err {
		*h = append(*h, val)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func FindKthSmallest(intHeap IntHeap, number int) int {
	if (intHeap.Len() < number) {
		errors.New("Heap size less than number")
	}

	sort.Sort(intHeap)

	for range number - 1 {
		intHeap.Pop()
	}
	
	return intHeap.Pop().(int)
}
