package intheap

type CustomHeap []int

func (h *CustomHeap) Len() int {
	return len(*h)
}

func (h *CustomHeap) Less(index1, index2 int) bool {
	if index1 < 0 || index1 >= len(*h) || index2 < 0 || index2 >= len(*h) {
		return false
	}

	return (*h)[index1] < (*h)[index2]
}

func (h *CustomHeap) Swap(index1, index2 int) {
	if index1 < 0 || index1 >= len(*h) || index2 < 0 || index2 >= len(*h) {
		return
	}

	(*h)[index1], (*h)[index2] = (*h)[index2], (*h)[index1]
}

func (h *CustomHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("IntHeap.Push: value isn't int")
	}

	*h = append(*h, value)
}

func (h *CustomHeap) Pop() interface{} {
	oldSlice := *h
	length := len(oldSlice)

	if length == 0 {
		return nil
	}

	lastElement := oldSlice[length-1]
	*h = oldSlice[0 : length-1]

	return lastElement
}
