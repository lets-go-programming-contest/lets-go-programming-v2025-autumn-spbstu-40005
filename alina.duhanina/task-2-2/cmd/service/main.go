package main

import (
	"container/heap"
	"errors"
	"fmt"

	"alina.duhanina/task-2-2/internal/intheap"
)

var (
	ErrEmptyDishes      = errors.New("Empty Dishes")
	ErrInvalidDishCount = errors.New("Invalid Dish Count")
	ErrInvalidK         = errors.New("Invalid K")
	ErrInvalidLenDishes = errors.New("Invalid LenDishes")
	ErrInvalidRating    = errors.New("Invalid Rating")
)

func ValidateInput(N int, dishes []int, k int) error {
	if len(dishes) == 0 {
		return ErrEmptyDishes
	}
	if N < 1 || N > 10000 {
		return ErrInvalidDishCount
	}
	if len(dishes) != N {
		return ErrInvalidLenDishes
	}
	if k < 1 || k > N {
		return ErrInvalidK
	}
	for _, rating := range dishes {
		if rating < -10000 || rating > 10000 {
			return ErrInvalidRating
		}
	}

	return nil
}

func FindKthPreference(dishes []int, k int) (int, error) {
	if err := ValidateInput(len(dishes), dishes, k); err != nil {
		return 0, err
	}

	h := &intheap.IntHeap{}
	heap.Init(h)
	for _, dish := range dishes {
		heap.Push(h, dish)
	}
	for i := 0; i < k-1; i++ {
		heap.Pop(h)
	}
	result := heap.Pop(h).(int)

	return result, nil
}

func main() {
	var N, k int

	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Printf("Invalid read: %v\n", err)
		return
	}

	dishes := make([]int, N)
	for i := 0; i < N; i++ {
		_, err := fmt.Scan(&dishes[i])
		if err != nil {
			fmt.Printf("Invalid read: %v\n", err)
			return
		}
	}
	_, err = fmt.Scan(&k)
	if err != nil {
		fmt.Printf("Invalid read: %v\n", err)
		return
	}

	result, err := FindKthPreference(dishes, k)
	if err != nil {
		fmt.Printf("Processing Error: %v\n", err)
		return
	}
	fmt.Println(result)
}
