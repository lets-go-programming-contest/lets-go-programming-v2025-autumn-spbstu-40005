package main

import (
	"container/heap"
	"errors"
	"fmt"

	"alina.duhanina/task-2-2/internal/intheap"
)

var (
	ErrEmptyDishes      = errors.New("empty dishes")
	ErrInvalidDishCount = errors.New("invalid dish count")
	ErrInvalidK         = errors.New("invalid k")
	ErrInvalidLenDishes = errors.New("invalid len dishes")
	ErrInvalidRating    = errors.New("invalid rating")
)

func ValidateInput(dishCount int, dishes []int, preferenceOrder int) error {
	if len(dishes) == 0 {
		return ErrEmptyDishes
	}

	if dishCount < 1 || dishCount > 10000 {
		return ErrInvalidDishCount
	}

	if len(dishes) != dishCount {
		return ErrInvalidLenDishes
	}

	if preferenceOrder < 1 || preferenceOrder > dishCount {
		return ErrInvalidK
	}

	for _, rating := range dishes {
		if rating < -10000 || rating > 10000 {
			return ErrInvalidRating
		}
	}

	return nil
}

func FindKthPreference(dishes []int, preferenceOrder int) (int, error) {
	if err := ValidateInput(len(dishes), dishes, preferenceOrder); err != nil {
		return 0, err
	}

	h := &intheap.IntHeap{}
	heap.Init(h)

	for _, dish := range dishes {
		heap.Push(h, dish)
	}
	for range preferenceOrder - 1 {
		heap.Pop(h)
	}

	result, ok := heap.Pop(h).(int)
	if !ok {
		return 0, errors.New("type assertion failed")
	}

	return result, nil
}

func main() {
	var dishCount, preferenceOrder int

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Printf("Invalid read: %v\n", err)
		return
	}

	dishes := make([]int, dishCount)
	for i := range dishCount {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			fmt.Printf("Invalid read: %v\n", err)

			return
		}
	}

	_, err = fmt.Scan(&preferenceOrder)
	if err != nil {
		fmt.Printf("Invalid read: %v\n", err)
		return
	}

	result, err := FindKthPreference(dishes, preferenceOrder)
	if err != nil {
		fmt.Printf("Processing Error: %v\n", err)
		return
	}

	fmt.Println(result)
}
