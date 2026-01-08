package heap

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

func FindKthPreferred(ratings []int, k int) int {
	h := &MinHeap{}

	for _, rating := range ratings {
		*h = append(*h, rating)

		if len(*h) > k {
			minIdx := 0
			for i := 1; i < len(*h); i++ {
				if (*h)[i] < (*h)[minIdx] {
					minIdx = i
				}
			}
			(*h)[minIdx] = (*h)[len(*h)-1]
			*h = (*h)[:len(*h)-1]
		}
	}

	min := (*h)[0]
	for i := 1; i < len(*h); i++ {
		if (*h)[i] < min {
			min = (*h)[i]
		}
	}
	return min
}
