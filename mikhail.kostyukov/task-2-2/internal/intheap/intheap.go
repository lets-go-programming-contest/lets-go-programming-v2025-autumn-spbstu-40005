package intheap

type IntHeap []int

func (heap *IntHeap) Len() int {
	return len(*heap)
}

func (heap *IntHeap) Less(first, second int) bool {
	if (first < 0) || (first >= heap.Len()) {
		panic("first index out of range")
	} else if (second < 0) || (second >= heap.Len()) {
		panic("second index out of range")
	}

	return (*heap)[first] > (*heap)[second]
}

func (heap *IntHeap) Swap(first, second int) {
	if (first < 0) || (first >= heap.Len()) {
		panic("first index out of range")
	} else if (second < 0) || (second >= heap.Len()) {
		panic("second index out of range")
	}

	(*heap)[first], (*heap)[second] = (*heap)[second], (*heap)[first]
}

func (heap *IntHeap) Push(elem any) {
	val, ok := elem.(int)
	if !ok {
		panic("pushed value is not int")
	}

	*heap = append(*heap, val)
}

func (heap *IntHeap) Pop() any {
	if heap.Len() == 0 {
		return nil
	}

	elem := (*heap)[len(*heap)-1]
	*heap = (*heap)[:len(*heap)-1]

	return elem
}
