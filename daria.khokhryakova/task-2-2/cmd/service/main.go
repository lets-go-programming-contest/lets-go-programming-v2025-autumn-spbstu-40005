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
	errFormat     = errors.New("invalid format")
	errOutOfRange = errors.New("invalid number")
)

func priorityDish(resHeap *intheap.IntHeap, num int) (int, error) {
	if num < minDishes || num > resHeap.Len() {
		return 0, fmt.Errorf("%w: %d out of range [%d, %d]", errOutOfRange, num, minDishes, resHeap.Len())
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
	if err != nil {
		fmt.Println("failed to read numberDishes:", err)

		return
	}

	if numberDishes < minDishes || numberDishes > maxDishes {

		fmt.Println("invalid input numberDishes")

		return
	}

	resHeap := &intheap.IntHeap{}
	heap.Init(resHeap)

	for range numberDishes {
		var rating int

		_, err := fmt.Scan(&rating)
		if err != nil || rating > maxRating || rating < minRating {
			fmt.Println("failed to read rating:", err)

			return
		}

		if rating > maxRating || rating < minRating {
			fmt.Println("invalid input rating")

			return
		}

		heap.Push(resHeap, rating)
	}

	var num int

	_, err = fmt.Scan(&num)
	if err != nil {
		fmt.Println("failed to read num:", err)

		return
	}

	if num > numberDishes || num < minDishes {
		fmt.Println("invalid input num")

		return
	}

	result, err := priorityDish(resHeap, num)
	if err != nil {
		fmt.Println("failed in the function priorityDish:", err)

		return
	}

	fmt.Println(result)
}
