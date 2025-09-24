package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/P3rCh1/task-2-2/internal/intheap"
)

var (
	errInputFail      = errors.New("input error")
	errInvalidRequest = errors.New("invalid request")
)

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil || dishesCount <= 0 {
		fmt.Println(errInputFail.Error())

		return
	}

	dishesSlice := make([]int, dishesCount)

	for index := range dishesCount {
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

	dishesHeap := new(intheap.IntHeap)
	*dishesHeap = intheap.IntHeap(dishesSlice[:need])
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
