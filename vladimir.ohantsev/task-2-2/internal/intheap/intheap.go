package intheap

type IntHeap []int

func (h *IntHeap) Push(val any) {
	intVal, ok := val.(int)
	if !ok {
		panic("invalid type pushed to IntHeap")
	}

	*h = append(*h, intVal)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		panic("heap underflow")
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
	h[i], h[j] = h[j], h[i]
}
