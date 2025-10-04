package main

import (
	"container/heap"
	"errors"
	"fmt"

	"gordey.shapkov/task-2-2/internal/intheap"
)

var errInvalidType = errors.New("invalid type")

func main() {
	var amount int
	if _, err := fmt.Scan(&amount); err != nil {
		return
	}

	dishes := &intheap.IntHeap{}

	for range amount {
		var pref int
		if _, err := fmt.Scan(&pref); err != nil {
			fmt.Println(errInvalidType)

			return
		}

		heap.Push(dishes, pref)
	}

	var number int
	if _, err := fmt.Scan(&number); err != nil {
		return
	}

	result, err := findDish(dishes, number)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func findDish(dishes *intheap.IntHeap, number int) (int, error) {
	for range dishes.Len() - number {
		heap.Pop(dishes)
	}

	x := heap.Pop(dishes)

	value, ok := x.(int)
	if !ok {
		return 0, errInvalidType
	}

	return value, nil
}
