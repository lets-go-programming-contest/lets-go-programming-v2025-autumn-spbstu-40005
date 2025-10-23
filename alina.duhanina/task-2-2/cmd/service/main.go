package main

import (
	"container/heap"
	"errors"
	"fmt"

	"alina.duhanina/task-2-2/internal/intheap"
)

var (
	ErrEmptyDishes      = errors.New("empty dishes")
	ErrInvalidK         = errors.New("invalid k")
	ErrInvalidLenDishes = errors.New("invalid len dishes")
	ErrInvalidRating    = errors.New("invalid rating")
	ErrAssertionFailed  = errors.New("type assertion failed")
)

func ValidateInput(dishCount int, dishes []int, preferenceOrder int) error {
	if len(dishes) == 0 {
		return ErrEmptyDishes
	}

	if len(dishes) != dishCount {
		return ErrInvalidLenDishes
	}

	if preferenceOrder < 1 || preferenceOrder > dishCount {
		return ErrInvalidK
	}

	return nil
}

func FindKthPreference(dishes []int, preferenceOrder int) (int, error) {
	heapInstance := &intheap.IntHeap{}
	heap.Init(heapInstance)

	for _, dish := range dishes {
		heap.Push(heapInstance, dish)
	}

	for range preferenceOrder - 1 {
		heap.Pop(heapInstance)
	}

	result, ok := heap.Pop(heapInstance).(int)
	if !ok {
		return 0, ErrAssertionFailed
	}

	return result, nil
}

func main() {
	var dishCount, preferenceOrder int

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Printf("Error reading number of dishes: %v\n", err)

		return
	}

	dishes := make([]int, dishCount)
	for i := range dishCount {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			fmt.Printf("Error reading rating for dish: %v\n", err)

			return
		}
	}

	_, err = fmt.Scan(&preferenceOrder)
	if err != nil {
		fmt.Printf("Error reading preference order: %v\n", err)

		return
	}

	if err := ValidateInput(len(dishes), dishes, preferenceOrder); err != nil {
		fmt.Printf("Input validation error: %v\n", err)

		return
	}

	result, err := FindKthPreference(dishes, preferenceOrder)
	if err != nil {
		fmt.Printf("Data processing error: %v\n", err)

		return
	}

	fmt.Println(result)
}
