package main

import (
	"container/heap"
	"errors"
	"fmt"

	"github.com/DariaKhokhryakova/task-2-2/internal/intheap"
)

const (
	minDishes = 1
	maxDishes = 10000
	minRating = -10000
	maxRating = 10000
)

var (
	errInput  = errors.New("invalid input")
	errFormat = errors.New("invalid format")
)

func priorityDish(resHeap *intheap.IntHeap, num int) (int, error) {
	if num < minDishes || num > resHeap.Len() {
		return 0, fmt.Errorf("invalid number %d out of range [%d, %d]", num, minDishes, resHeap.Len())
	}

	for range num - 1 {
		heap.Pop(resHeap)
	}

	resPop := heap.Pop(resHeap)
	resPriority, ok := resPop.(int)

	if !ok {
		return 0, errFormat
	}

	return resPriority, nil
}

func main() {
	var numberDishes int

	_, err := fmt.Scan(&numberDishes)
	if err != nil || numberDishes < minDishes || numberDishes > maxDishes {
		fmt.Println(errInput.Error()+":", err)

		return
	}

	resHeap := &intheap.IntHeap{}
	heap.Init(resHeap)

	for range numberDishes {
		var rating int

		_, err := fmt.Scan(&rating)
		if err != nil || rating > maxRating || rating < minRating {
			fmt.Println(errInput.Error())

			return
		}

		heap.Push(resHeap, rating)
	}

	var num int

	_, err = fmt.Scan(&num)
	if err != nil || num > numberDishes || num < minDishes {
		fmt.Println(errInput.Error())

		return
	}

	result, err := priorityDish(resHeap, num)
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Println(result)
}
