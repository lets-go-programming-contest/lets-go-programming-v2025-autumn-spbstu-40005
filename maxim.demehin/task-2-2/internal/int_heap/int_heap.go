package int_heap

import (
	"container/heap"
)

type IntHeap []int

func (h *IntHeap) Push(val any) {
	*h = append(*h, val.(int))
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		panic("error: heap underflow")
	}

	orig := *h
	len := len(orig)
	toreturn := orig[len-1]
	*h = orig[0 : len-1]
	return toreturn
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h IntHeap) Swap(i, j int) {
	tmp := h[i]
	h[i] = h[j]
	h[j] = tmp
}

func (h *IntHeap) GetNth(n int) int {
	temp := make(IntHeap, h.Len())
	copy(temp, *h)
	heap.Init(&temp)

	var res int
	for range n {
		res = heap.Pop(&temp).(int)
	}

	return res
}
