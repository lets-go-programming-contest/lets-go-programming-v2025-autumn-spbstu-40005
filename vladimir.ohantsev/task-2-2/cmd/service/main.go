package main

import (
	"container/heap"
	"fmt"

	"github.com/P3rCh1/task-2-2/internal/intheap"
)

func getKDish(dishes []int, k int) int {
	window := new(intheap.IntHeap)
	*window = intheap.IntHeap(dishes[:k])
	heap.Init(window)

	dishes = dishes[k:]
	for _, cost := range dishes {
		top, err := window.Top()
		if err != nil {
			panic(fmt.Sprintf("top heap: %s", err))
		}

		if cost > top {
			err := window.ReplaceTop(cost)
			if err != nil {
				panic(fmt.Sprintf("replace top heap: %s", err))
			}
		}
	}

	top, err := window.Top()
	if err != nil {
		panic(fmt.Sprintf("top heap: %s", err))
	}

	return top
}

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Printf("failed to scan dishes count: %s\n", err)

		return
	}

	if dishesCount <= 0 {
		fmt.Printf("non-positive dishes count")

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
	if _, err := fmt.Scan(&need); err != nil {
		fmt.Printf("failed to scan needed element: %s\n", err)

		return
	}

	if need <= 0 {
		fmt.Printf("non-positive needed item")

		return
	}

	if need > dishesCount {
		fmt.Println("dishes count should be less then needed element")

		return
	}

	fmt.Println(getKDish(dishes, need))
}
