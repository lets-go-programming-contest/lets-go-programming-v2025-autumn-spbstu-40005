package minheap

type MinHeap []int

func (heap *MinHeap) Len() int {
	return len(*heap)
}

func (heap *MinHeap) Less(i, j int) bool {
	return (*heap)[i] < (*heap)[j]
}

func (heap *MinHeap) Swap(i, j int) {
	(*heap)[i], (*heap)[j] = (*heap)[j], (*heap)[i]
}

func (heap *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("invalid argument")
	}

	*heap = append(*heap, value)
}

func (heap *MinHeap) Pop() interface{} {
	if len(*heap) == 0 {
		return nil
	}

	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

func (heap *MinHeap) GetMin() int {
	if len(*heap) == 0 {
		panic("heap is empty")
	}

	return (*heap)[0]
}
