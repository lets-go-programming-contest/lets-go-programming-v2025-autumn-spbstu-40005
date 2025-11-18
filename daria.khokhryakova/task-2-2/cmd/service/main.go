package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/DariaKhokhryakova/task-2-2/internal/intheap"
)

var (
	errFormat     = errors.New("invalid format")
	errOutOfRange = errors.New("invalid number")
)

func priorityDish(resHeap *intheap.IntHeap, num int) (int, error) {
	for range num - 1 {
		heap.Pop(resHeap)
	}

	resPop := heap.Pop(resHeap)
	if resPop == nil {
		return 0, errOutOfRange
	}

	resPriority, ok := resPop.(int)
	if !ok {
		return 0, errFormat
	}

	return resPriority, nil
}

func main() {
	var numberDishes int

	_, err := fmt.Scan(&numberDishes)
	if err != nil {
		fmt.Println("error in the numberDishes parameter:", err)

		return
	}

	resHeap := &intheap.IntHeap{}
	heap.Init(resHeap)

	for range numberDishes {
		var rating int

		_, err := fmt.Scan(&rating)
		if err != nil {
			fmt.Println("error in the rating parameter:", err)

			return
		}

		heap.Push(resHeap, rating)
	}

	var num int

	_, err = fmt.Scan(&num)
	if err != nil {
		fmt.Println("error in the num parameter:", err)

		return
	}

	result, err := priorityDish(resHeap, num)
	if err != nil {
		fmt.Println("failed in the function priorityDish:", err)

		return
	}

	fmt.Println(result)
}
