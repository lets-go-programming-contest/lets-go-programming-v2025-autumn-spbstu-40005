package intheap

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Swap(i, j int) {
	if i >= h.Len() || j >= h.Len() || i < 0 || j < 0 {
		panic("index out of range in Swap")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Less(i, j int) bool {
	if i >= h.Len() || j >= h.Len() || i < 0 || j < 0 {
		panic("index out of range in Less")
	}
	
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("Push expects int")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	if h.Len() == 0 {
		return nil
	}

	orig := *h
	n := len(orig)
	x := orig[n-1]
	*h = orig[0 : n-1]

	return x
}
