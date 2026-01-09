package main

import (
	"container/heap"
	"errors"
	"fmt"

	"aleksey.kurbyko/task-2-2/internal/dishheap"
)

const (
	DishCountMin = 1
	DishCountMax = 10000

	RatingMin = -10000
	RatingMax = 10000
)

var (
	ErrIncorrectDishCount = errors.New("incorrect amount of dishes")
	ErrIncorrectRating    = errors.New("incorrect rating for the dish")
	ErrIncorrectK         = errors.New("incorrect k")
	ErrExpectedInt        = errors.New("expected int")
	ErrEmptyHeap          = errors.New("empty heap")
)

func readInt() (int, error) {
	var value int

	if _, err := fmt.Scan(&value); err != nil {
		return 0, fmt.Errorf("scan int: %w", err)
	}

	return value, nil
}

func readDishCount() (int, error) {
	dishCount, err := readInt()
	if err != nil {
		return 0, ErrIncorrectDishCount
	}

	if dishCount < DishCountMin || dishCount > DishCountMax {
		return 0, ErrIncorrectDishCount
	}

	return dishCount, nil
}

func readPreferredIndex(dishCount int) (int, error) {
	preferredIndex, err := readInt()
	if err != nil {
		return 0, ErrIncorrectK
	}

	if preferredIndex < 1 || preferredIndex > dishCount {
		return 0, ErrIncorrectK
	}

	return preferredIndex, nil
}

func buildRatingsHeap(dishCount int) (*dishheap.DishHeap, error) {
	ratingsHeap := &dishheap.DishHeap{}
	heap.Init(ratingsHeap)

	for range dishCount {
		rating, err := readInt()
		if err != nil {
			return nil, ErrIncorrectRating
		}

		if rating < RatingMin || rating > RatingMax {
			return nil, ErrIncorrectRating
		}

		heap.Push(ratingsHeap, rating)
	}

	return ratingsHeap, nil
}

func getPreferredRating(ratingsHeap *dishheap.DishHeap, preferredIndex int) (int, error) {
	for range preferredIndex - 1 {
		heap.Pop(ratingsHeap)
	}

	popped := heap.Pop(ratingsHeap)
	if popped == nil {
		return 0, ErrEmptyHeap
	}

	result, ok := popped.(int)
	if !ok {
		return 0, ErrExpectedInt
	}

	return result, nil
}

func run() error {
	dishCount, err := readDishCount()
	if err != nil {
		return err
	}

	ratingsHeap, err := buildRatingsHeap(dishCount)
	if err != nil {
		return err
	}

	preferredIndex, err := readPreferredIndex(dishCount)
	if err != nil {
		return err
	}

	result, err := getPreferredRating(ratingsHeap, preferredIndex)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func main() {
	if err := run(); err != nil {
		return
	}
}
