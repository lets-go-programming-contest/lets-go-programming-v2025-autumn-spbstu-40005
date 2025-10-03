package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/P3rCh1/task-2-2/internal/intheap"
)

var errInvalidNeed = errors.New("dishes count should be less then needed element")

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil || dishesCount <= 0 {
		fmt.Printf("failed to dishes count: %s\n", err)

		return
	}

	dishes := make([]int, dishesCount)

	for index := range dishesCount {
		var cost int
		if _, err := fmt.Scan(&cost); err != nil {
			fmt.Printf("failed to scan cost: %s\n", err)

			return
		}

		dishes[index] = cost
	}

	var need int
	if _, err := fmt.Scan(&need); err != nil || need <= 0 {
		fmt.Printf("failed to scan needed element: %s\n", err)

		return
	}

	if need > dishesCount {
		fmt.Println(errInvalidNeed.Error())

		return
	}

	window := new(intheap.IntHeap)
	*window = intheap.IntHeap(dishes[:need])
	heap.Init(window)

	dishes = dishes[need:]
	for _, cost := range dishes {
		if cost > window.Top() {
			window.ReplaceTop(cost)
		}
	}

	fmt.Println(window.Top())
}
