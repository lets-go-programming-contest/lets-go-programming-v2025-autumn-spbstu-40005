package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	errInputFail      = errors.New("input error")
	errInvalidRequest = errors.New("invalid request")
)

type intHeap []int

func (h *intHeap) Push(x any) {
	if x, ok := x.(int); ok {
		*h = append(*h, x)
	}
}

func (h *intHeap) Pop() any {
	val := (*h)[0]
	*h = (*h)[1:]

	return val
}

func (h *intHeap) Len() int           { return len(*h) }
func (h *intHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *intHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil || dishesCount <= 0 {
		fmt.Println(errInputFail.Error())

		return
	}

	dishesSlice := make([]int, dishesCount)

	for index := 0; index < dishesCount; index++ {
		var cost int
		if _, err := fmt.Scan(&cost); err != nil {
			fmt.Println(errInputFail.Error())

			return
		}

		dishesSlice[index] = cost
	}

	var need int
	if _, err := fmt.Scan(&need); err != nil || need <= 0 {
		fmt.Println(errInputFail.Error())

		return
	}

	if need > dishesCount {
		fmt.Println(errInvalidRequest.Error())

		return
	}

	dishesHeap := &intHeap{}
	*dishesHeap = intHeap(dishesSlice[:need])
	heap.Init(dishesHeap)

	dishesSlice = dishesSlice[need:]
	for _, cost := range dishesSlice {
		if cost > (*dishesHeap)[0] {
			(*dishesHeap)[0] = cost
			heap.Fix(dishesHeap, 0)
		}
	}

	fmt.Println((*dishesHeap)[0])
}
