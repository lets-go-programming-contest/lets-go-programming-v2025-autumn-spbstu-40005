package queue

type Queue []int

func (q *Queue) Len() int {
	return len(*q)
}

func (q *Queue) Less(i, j int) bool {
	if i >= q.Len() || j >= q.Len() || i < 0 || j < 0 {
		panic("indexes out of range")
	}

	return (*q)[i] > (*q)[j]
}

func (q *Queue) Swap(i, j int) {
	if i >= q.Len() || j >= q.Len() || i < 0 || j < 0 {
		panic("swapping out of range")
	}

	(*q)[i], (*q)[j] = (*q)[j], (*q)[i]
}

func (q *Queue) Push(newElement any) {
	element, ok := newElement.(int)
	if !ok {
		panic("queue only accepts integer values")
	}

	*q = append(*q, element)
}

func (q *Queue) Pop() any {
	if q.Len() == 0 {
		return nil
	}

	old := *q
	n := len(old)
	item := old[n-1]

	*q = old[0 : n-1]

	return item
}
