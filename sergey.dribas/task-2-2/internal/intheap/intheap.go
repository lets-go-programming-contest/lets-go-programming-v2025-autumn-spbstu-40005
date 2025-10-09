package intheap

import (
	"fmt"
	"sort"
)

var (
	ErrHeap = fmt.Errorf("heap size less than number")
	ErrType = fmt.Errorf("unexpected type from heap pop")
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

func FindKthSmallest(intHeap IntHeap, number int) (int, error) {
	if intHeap.Len() < number {
		return 0, ErrHeap
	}

	sort.Sort(intHeap)

	for range number - 1 {
		intHeap.Pop()
	}

	result := intHeap.Pop()
	if val, ok := result.(int); ok {
		return val, nil
	}

	return 0, ErrType
}
