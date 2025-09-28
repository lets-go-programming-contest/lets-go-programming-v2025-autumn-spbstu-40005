package IntHeap

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	temp := h[i]
	h[i] = h[j]
	h[j] = temp
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("Invalid type")
	}
	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		panic("Underflow")
	}
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
