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
	val := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]

	return val
}

func (h *intHeap) Len() int           { return len(*h) }
func (h *intHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *intHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func main() {
	var dishes int
	if _, err := fmt.Scan(&dishes); err != nil || dishes <= 0 {
		fmt.Println(errInputFail.Error())

		return
	}

	dishesHeap := &intHeap{}

	for range dishes {
		var cost int
		if _, err := fmt.Scan(&cost); err != nil {
			fmt.Println(errInputFail.Error())

			return
		}

		heap.Push(dishesHeap, cost)
	}

	var need int
	if _, err := fmt.Scan(&need); err != nil || need <= 0 {
		fmt.Println(errInputFail.Error())

		return
	}

	if need > dishesHeap.Len() {
		fmt.Println(errInvalidRequest.Error())

		return
	}

	for range need - 1 {
		heap.Pop(dishesHeap)
	}

	res, _ := heap.Pop(dishesHeap).(int)
	fmt.Println(res)
}
