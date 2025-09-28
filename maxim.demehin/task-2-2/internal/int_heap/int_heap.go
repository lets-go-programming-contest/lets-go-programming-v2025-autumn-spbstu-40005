package int_heap

type IntHeap []int

func (h *IntHeap) Push(val any) {
	*h = append(*h, val.(int))
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		panic("error: heap underflow")
	}

	val := (*h)[0]
	*h = (*h)[1:]

	return val
}

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	tmp := h[i]
	h[i] = h[j]
	h[j] = tmp
}
